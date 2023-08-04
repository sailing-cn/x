package device

import (
	"fmt"
	"github.com/golang/glog"
	"sailing.cn/v2/convert"
	"sailing.cn/v2/iot"
	mqtt2 "sailing.cn/v2/mqtt"
	"sailing.cn/v2/utils/timestamp"
	"time"
)

type Device interface {
	BaseDevice
	Gateway
	SendMessage(message iot.Message) bool
	ReportProperties(properties iot.DeviceProperties) bool
	BatchReportSubDevicesProperties(service iot.DevicesService) bool
	QueryDeviceShadow(query iot.DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler)
	UploadFile(filename string) bool
	DownloadFile(filename string) bool
	ReportDeviceInfo(version *iot.VersionInfo)
	ReportLogs(logs []iot.DeviceLog) bool
}

type iotDevice struct {
	base    baseIotDevice
	gateway baseGateway
}

func (device *iotDevice) Init() bool {
	return device.base.Init()
}

func (device *iotDevice) DisConnect() {
	device.base.DisConnect()
}

func (device *iotDevice) IsConnected() bool {
	return device.base.IsConnected()
}

func (device *iotDevice) AddMessageHandler(handler MessageHandler) {
	device.base.AddMessageHandler(handler)
}
func (device *iotDevice) AddCommandHandler(handler CommandHandler) {
	device.base.AddCommandHandler(handler)
}
func (device *iotDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	device.base.AddPropertiesSetHandler(handler)
}
func (device *iotDevice) SetVersionReporter(handler VersionReporter) {
	device.base.SetVersionReporter(handler)
}

func (device *iotDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.base.SetDeviceUpgradeHandler(handler)
}

func (device *iotDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.base.SetPropertyQueryHandler(handler)
}

func (device *iotDevice) ReportLogs(logs []iot.DeviceLog) bool {
	var services []iot.ReportDeviceLogServiceEvent

	for _, logEntry := range logs {
		service := iot.ReportDeviceLogServiceEvent{
			BaseServiceEvent: iot.BaseServiceEvent{
				ServiceId: "$log",
				EventType: "log_report",
				EventTime: timestamp.Timestamp(),
			},
			Paras: logEntry,
		}

		services = append(services, service)
	}

	request := iot.ReportDeviceLogRequest{
		Services: services,
	}

	fmt.Println(convert.ToString(request))

	topic := mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id)

	token := device.base.Client.Publish(topic, 1, false, convert.ToString(request))

	if token.Wait() && token.Error() != nil {
		glog.Errorf("device %s report log failed", device.base.Id)
		return false
	} else {
		return true
	}
}

func (device *iotDevice) SendMessage(message iot.Message) bool {
	messageData := convert.ToString(message)
	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.MessageUpTopic, device.base.Id), device.base.qos, false, messageData); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s send message failed", device.base.Id)
		return false
	}
	return true
}

func (device *iotDevice) ReportProperties(properties iot.DeviceProperties) bool {
	propertiesData := convert.ToString(properties)
	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.PropertiesUpTopic, device.base.Id), device.base.qos, false, propertiesData); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s report properties failed", device.base.Id)
		return false
	}
	return true
}
func (device *iotDevice) BatchReportSubDevicesProperties(service iot.DevicesService) bool {

	subDeviceCounts := len(service.Devices)

	batchReportSubDeviceProperties := 0
	if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
		batchReportSubDeviceProperties = subDeviceCounts / device.base.batchSubDeviceSize
	} else {
		batchReportSubDeviceProperties = subDeviceCounts/device.base.batchSubDeviceSize + 1
	}

	for i := 0; i < batchReportSubDeviceProperties; i++ {
		begin := i * device.base.batchSubDeviceSize
		end := (i + 1) * device.base.batchSubDeviceSize
		if end > subDeviceCounts {
			end = subDeviceCounts
		}

		sds := iot.DevicesService{
			Devices: service.Devices[begin:end],
		}

		if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.GatewayBatchReportSubDeviceTopic, device.base.Id), device.base.qos, false, convert.ToString(sds)); token.Wait() && token.Error() != nil {
			glog.Warningf("device %s batch report sub device properties failed", device.base.Id)
			return false
		}
	}

	return true
}

func (device *iotDevice) QueryDeviceShadow(query iot.DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) {
	device.base.propertiesQueryResponseHandler = handler
	requestId := timestamp.TimestampString()
	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceShadowQueryRequestTopic, device.base.Id)+requestId, device.base.qos, false, convert.ToString(query)); token.Wait() && token.Error() != nil {
		glog.Warningf("device %s query device shadow data failed,request id = %s", device.base.Id, requestId)
	}
}

func (device *iotDevice) UploadFile(filename string) bool {
	// 构造获取文件上传URL的请求
	requestParas := iot.FileRequestServiceEventParas{
		FileName: filename,
	}

	serviceEvent := iot.FileRequestServiceEvent{
		Paras: requestParas,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventTime = timestamp.Timestamp()
	serviceEvent.EventType = "get_upload_url"

	var services []iot.FileRequestServiceEvent
	services = append(services, serviceEvent)
	request := iot.FileRequest{
		Services: services,
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request)); token.Wait() && token.Error() != nil {
		glog.Warningf("publish file upload request url failed")
		return false
	}
	glog.Info("publish file upload request url success")

	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			_, ok := device.base.fileUrls[filename+iot.FileActionUpload]
			if ok {
				glog.Infof("platform send file upload url success")
				goto BreakPoint
			}

		}
	}
BreakPoint:

	if len(device.base.fileUrls[filename+iot.FileActionUpload]) == 0 {
		glog.Errorf("get file upload url failed")
		return false
	}
	glog.Infof("file upload url is %s", device.base.fileUrls[filename+iot.FileActionUpload])

	//filename = smartFileName(filename)
	uploadFlag := CreateHttpClient().UploadFile(filename, device.base.fileUrls[filename+iot.FileActionUpload])
	if !uploadFlag {
		glog.Errorf("upload file failed")
		return false
	}

	response := CreateFileUploadDownLoadResultResponse(filename, iot.FileActionUpload, uploadFlag)

	token := device.base.Client.Publish(mqtt2.FormatTopic(iot.PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, convert.ToString(response))
	if token.Wait() && token.Error() != nil {
		glog.Error("report file upload file result failed")
		return false
	}

	return true
}

func (device *iotDevice) DownloadFile(filename string) bool {
	// 构造获取文件上传URL的请求
	requestParas := iot.FileRequestServiceEventParas{
		FileName: filename,
	}

	serviceEvent := iot.FileRequestServiceEvent{
		Paras: requestParas,
	}
	serviceEvent.ServiceId = "$file_manager"
	serviceEvent.EventTime = timestamp.Timestamp()
	serviceEvent.EventType = "get_download_url"

	var services []iot.FileRequestServiceEvent
	services = append(services, serviceEvent)
	request := iot.FileRequest{
		Services: services,
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request)); token.Wait() && token.Error() != nil {
		glog.Warningf("publish file download request url failed")
		return false
	}

	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			_, ok := device.base.fileUrls[filename+iot.FileActionDownload]
			if ok {
				glog.Infof("platform send file upload url success")
				goto BreakPoint
			}

		}
	}
BreakPoint:

	if len(device.base.fileUrls[filename+iot.FileActionDownload]) == 0 {
		glog.Errorf("get file download url failed")
		return false
	}

	downloadFlag := CreateHttpClient().DownloadFile(filename, device.base.fileUrls[filename+iot.FileActionDownload])
	if !downloadFlag {
		glog.Errorf("down load file { %s } failed", filename)
		return false
	}

	response := CreateFileUploadDownLoadResultResponse(filename, iot.FileActionDownload, downloadFlag)

	token := device.base.Client.Publish(mqtt2.FormatTopic(iot.PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, convert.ToString(response))
	if token.Wait() && token.Error() != nil {
		glog.Error("report file upload file result failed")
		return false
	}

	return true
}

func (device *iotDevice) ReportDeviceInfo(version *iot.VersionInfo) {
	event := iot.ReportDeviceInfoServiceEvent{
		BaseServiceEvent: iot.BaseServiceEvent{
			ServiceId: "$device_info",
			EventType: "device_info_report",
			EventTime: timestamp.Timestamp(),
		},
		Paras: iot.ReportDeviceInfoEventParas{
			DeviceSdkVersion: iot.SdkInfo()["sdk-version"],
			SwVersion:        version.SoftwareVersion,
			FwVersion:        version.FirmwareVersion,
		},
	}

	request := iot.ReportDeviceInfoRequest{
		DeviceId: device.base.Id,
		Services: []iot.ReportDeviceInfoServiceEvent{event},
	}

	device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request))
}

func (device *iotDevice) SetSubDevicesAddHandler(handler SubDevicesAddHandler) {
	device.base.subDevicesAddHandler = handler
}

func (device *iotDevice) SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler) {
	device.base.subDevicesDeleteHandler = handler
}

func (device *iotDevice) SetDeviceStatusLogCollector(collector DeviceStatusLogCollector) {
	device.base.SetDeviceStatusLogCollector(collector)
}

func (device *iotDevice) SetDevicePropertyLogCollector(collector DevicePropertyLogCollector) {
	device.base.SetDevicePropertyLogCollector(collector)
}

func (device *iotDevice) SetDeviceMessageLogCollector(collector DeviceMessageLogCollector) {
	device.base.SetDeviceMessageLogCollector(collector)
}

func (device *iotDevice) SetDeviceCommandLogCollector(collector DeviceCommandLogCollector) {
	device.base.SetDeviceCommandLogCollector(collector)
}

func (device *iotDevice) UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) bool {
	glog.Infof("begin to update sub-devices status")

	subDeviceCounts := len(subDevicesStatus.DeviceStatuses)

	batchUpdateSubDeviceState := 0
	if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
		batchUpdateSubDeviceState = subDeviceCounts / device.base.batchSubDeviceSize
	} else {
		batchUpdateSubDeviceState = subDeviceCounts/device.base.batchSubDeviceSize + 1
	}

	for i := 0; i < batchUpdateSubDeviceState; i++ {
		begin := i * device.base.batchSubDeviceSize
		end := (i + 1) * device.base.batchSubDeviceSize
		if end > subDeviceCounts {
			end = subDeviceCounts
		}

		sds := iot.SubDevicesStatus{
			DeviceStatuses: subDevicesStatus.DeviceStatuses[begin:end],
		}

		requestEventService := iot.ServiceEvent{
			ServiceId: "$sub_device_manager",
			EventType: "sub_device_update_status",
			EventTime: timestamp.Timestamp(),
			Paras:     sds,
		}

		request := iot.EventData{
			DeviceId: device.base.Id,
			Services: []iot.ServiceEvent{requestEventService},
		}

		if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request)); token.Wait() && token.Error() != nil {
			glog.Warningf("gateway %s update sub devices status failed", device.base.Id)
			return false
		}
	}

	glog.Info("gateway  update sub devices status failed", device.base.Id)
	return true
}

func (device *iotDevice) DeleteSubDevices(deviceIds []string) bool {
	glog.Infof("begin to delete sub-devices %s", deviceIds)

	subDevices := struct {
		Devices []string `json:"devices"`
	}{
		Devices: deviceIds,
	}

	requestEventService := iot.ServiceEvent{
		ServiceId: "$sub_device_manager",
		EventType: "delete_sub_device_request",
		EventTime: timestamp.Timestamp(),
		Paras:     subDevices,
	}

	request := iot.EventData{
		DeviceId: device.base.Id,
		Services: []iot.ServiceEvent{requestEventService},
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request)); token.Wait() && token.Error() != nil {
		glog.Warningf("gateway %s delete sub devices request send failed", device.base.Id)
		return false
	}

	glog.Warningf("gateway %s delete sub devices request send success", device.base.Id)
	return true
}

func (device *iotDevice) AddSubDevices(deviceInfos []iot.SubDeviceInfo) bool {
	devices := struct {
		Devices []iot.SubDeviceInfo `json:"devices"`
	}{
		Devices: deviceInfos,
	}

	requestEventService := iot.ServiceEvent{
		ServiceId: "$sub_device_manager",
		EventType: "add_sub_device_request",
		EventTime: timestamp.Timestamp(),
		Paras:     devices,
	}

	request := iot.EventData{
		DeviceId: device.base.Id,
		Services: []iot.ServiceEvent{requestEventService},
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(request)); token.Wait() && token.Error() != nil {
		glog.Warningf("gateway %s add sub devices request send failed", device.base.Id)
		return false
	}

	glog.Warningf("gateway %s add sub devices request send success", device.base.Id)
	return true
}

func (device *iotDevice) SyncAllVersionSubDevices() {
	dataEntry := iot.ServiceEvent{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: timestamp.Timestamp(),
		Paras: struct {
		}{},
	}

	var dataEntries []iot.ServiceEvent
	dataEntries = append(dataEntries, dataEntry)

	data := iot.EventData{
		Services: dataEntries,
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("send sub device sync request failed")
	}
}

func (device *iotDevice) SyncSubDevices(version int) {
	syncParas := struct {
		Version int `json:"version"`
	}{
		Version: version,
	}

	dataEntry := iot.ServiceEvent{
		ServiceId: "$sub_device_manager",
		EventType: "sub_device_sync_request",
		EventTime: timestamp.Timestamp(),
		Paras:     syncParas,
	}

	var dataEntries []iot.ServiceEvent
	dataEntries = append(dataEntries, dataEntry)

	data := iot.EventData{
		Services: dataEntries,
	}

	if token := device.base.Client.Publish(mqtt2.FormatTopic(iot.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, convert.ToString(data)); token.Wait() && token.Error() != nil {
		glog.Errorf("send sync sub device request failed")
	}
}

func CreateIotDevice(id, password, servers string) Device {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      0,
		AuthType: AUTH_TYPE_PASSWORD,
	}

	return CreateIotDeviceWitConfig(config)
}
func CreateIotDevice1(id, password, servers, serverCaPath, certFilePath, certKeyFilePath string) Device {
	config := DeviceConfig{
		Id:              id,
		Password:        password,
		Servers:         servers,
		Qos:             0,
		AuthType:        AUTH_TYPE_X509,
		CertKeyFilePath: certKeyFilePath,
		CertFilePath:    certFilePath,
		ServerCaPath:    serverCaPath,
	}

	return CreateIotDeviceWitConfig(config)
}

func CreateIotDeviceWithQos(id, password, servers string, qos byte) Device {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      qos,
	}

	return CreateIotDeviceWitConfig(config)
}

func CreateIotDeviceWitConfig(config DeviceConfig) Device {
	device := baseIotDevice{}
	device.Id = config.Id
	device.Password = config.Password
	device.Servers = config.Servers
	device.messageHandlers = []MessageHandler{}
	device.commandHandlers = []CommandHandler{}

	device.fileUrls = map[string]string{}

	device.qos = config.Qos
	device.batchSubDeviceSize = 100
	device.AuthType = config.AuthType
	device.ServerCaPath = config.ServerCaPath
	device.CertFilePath = config.CertFilePath
	device.CertKeyFilePath = config.CertKeyFilePath

	result := &iotDevice{
		base: device,
	}
	return result
}
