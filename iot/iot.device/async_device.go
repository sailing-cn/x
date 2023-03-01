package iot_device

import (
	"fmt"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
	"sailing.cn/iot"
	"sailing.cn/iot/iot.device/topics"
	"time"
)

type AsyncDevice interface {
	BaseDevice
	AsyncGateway
	SendMessage(message iot.Message) AsyncResult
	ReportProperties(properties iot.DeviceProperties) AsyncResult
	BatchReportSubDevicesProperties(service iot.DevicesService) AsyncResult
	QueryDeviceShadow(query iot.DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) AsyncResult
	UploadFile(filename string) AsyncResult
	DownloadFile(filename string) AsyncResult
	ReportDeviceInfo(swVersion, fwVersion string) AsyncResult
	ReportLogs(logs []iot.DeviceLog) AsyncResult
}

func CreateAsyncIotDevice(id, password, servers string) *asyncDevice {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      0,
	}

	return CreateAsyncIotDeviceWitConfig(config)
}

func CreateAsyncIotDeviceWithQos(id, password, servers string, qos byte) *asyncDevice {
	config := DeviceConfig{
		Id:       id,
		Password: password,
		Servers:  servers,
		Qos:      qos,
	}

	return CreateAsyncIotDeviceWitConfig(config)
}

func CreateAsyncIotDeviceWitConfig(config DeviceConfig) *asyncDevice {
	device := baseIotDevice{}
	device.Id = config.Id
	device.Password = config.Password
	device.Servers = config.Servers
	device.messageHandlers = []MessageHandler{}
	device.commandHandlers = []CommandHandler{}

	device.fileUrls = map[string]string{}

	device.qos = config.Qos
	device.batchSubDeviceSize = config.BatchSubDeviceSize

	result := &asyncDevice{
		base: device,
	}
	return result
}

type asyncDevice struct {
	base baseIotDevice
}

func (device *asyncDevice) Init() bool {
	return device.base.Init()
}

func (device *asyncDevice) DisConnect() {
	device.base.DisConnect()
}

func (device *asyncDevice) IsConnected() bool {
	return device.base.IsConnected()
}

func (device *asyncDevice) AddMessageHandler(handler MessageHandler) {
	device.base.AddMessageHandler(handler)
}
func (device *asyncDevice) AddCommandHandler(handler CommandHandler) {
	device.base.AddCommandHandler(handler)
}
func (device *asyncDevice) AddPropertiesSetHandler(handler DevicePropertiesSetHandler) {
	device.base.AddPropertiesSetHandler(handler)
}
func (device *asyncDevice) SetPropertyQueryHandler(handler DevicePropertyQueryHandler) {
	device.base.SetPropertyQueryHandler(handler)
}

func (device *asyncDevice) SetVersionReporter(handler VersionReporter) {
	device.base.SetVersionReporter(handler)
}

func (device *asyncDevice) SetDeviceUpgradeHandler(handler DeviceUpgradeHandler) {
	device.base.SetDeviceUpgradeHandler(handler)
}

func (device *asyncDevice) SetDeviceStatusLogCollector(collector DeviceStatusLogCollector) {
	device.base.SetDeviceStatusLogCollector(collector)
}

func (device *asyncDevice) SetDevicePropertyLogCollector(collector DevicePropertyLogCollector) {
	device.base.SetDevicePropertyLogCollector(collector)
}

func (device *asyncDevice) SetDeviceMessageLogCollector(collector DeviceMessageLogCollector) {
	device.base.SetDeviceMessageLogCollector(collector)
}

func (device *asyncDevice) SetDeviceCommandLogCollector(collector DeviceCommandLogCollector) {
	device.base.SetDeviceCommandLogCollector(collector)
}

func (device *asyncDevice) SendMessage(message iot.Message) AsyncResult {
	asyncResult := NewBooleanAsyncResult()
	go func() {
		glog.Info("begin async send message")

		messageData := iot.Interface2JsonString(message)
		topic := iot.FormatTopic(topics.MessageUpTopic, device.base.Id)
		glog.Infof("async send message topic is %s", topic)
		token := device.base.Client.Publish(topic, device.base.qos, false, messageData)
		if token.Wait() && token.Error() != nil {
			glog.Warning("async send message failed")
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) ReportProperties(properties iot.DeviceProperties) AsyncResult {
	asyncResult := NewBooleanAsyncResult()
	go func() {
		glog.Info("begin to report properties")
		propertiesData := iot.Interface2JsonString(properties)
		if token := device.base.Client.Publish(iot.FormatTopic(topics.PropertiesUpTopic, device.base.Id), device.base.qos, false, propertiesData); token.Wait() && token.Error() != nil {
			glog.Warningf("device %s async report properties failed", device.base.Id)
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) BatchReportSubDevicesProperties(service iot.DevicesService) AsyncResult {
	asyncResult := NewBooleanAsyncResult()

	go func() {
		glog.Info("begin async batch report sub devices properties")
		subDeviceCounts := len(service.Devices)
		batchReportSubDeviceProperties := 0
		if subDeviceCounts%device.base.batchSubDeviceSize == 0 {
			batchReportSubDeviceProperties = subDeviceCounts / device.base.batchSubDeviceSize
		} else {
			batchReportSubDeviceProperties = subDeviceCounts/device.base.batchSubDeviceSize + 1
		}

		loopResult := true
		for i := 0; i < batchReportSubDeviceProperties; i++ {
			begin := i * device.base.batchSubDeviceSize
			end := (i + 1) * device.base.batchSubDeviceSize
			if end > subDeviceCounts {
				end = subDeviceCounts
			}

			sds := iot.DevicesService{
				Devices: service.Devices[begin:end],
			}

			if token := device.base.Client.Publish(iot.FormatTopic(topics.GatewayBatchReportSubDeviceTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(sds)); token.Wait() && token.Error() != nil {
				glog.Warningf("device %s batch report sub device properties failed", device.base.Id)
				loopResult = false
				asyncResult.completeError(token.Error())
				break
			}
		}

		if loopResult {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) QueryDeviceShadow(query iot.DevicePropertyQueryRequest, handler DevicePropertyQueryResponseHandler) AsyncResult {
	device.base.propertiesQueryResponseHandler = handler
	asyncResult := NewBooleanAsyncResult()

	go func() {
		requestId := uuid.NewV4()
		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceShadowQueryRequestTopic, device.base.Id)+requestId.String(), device.base.qos, false, iot.Interface2JsonString(query)); token.Wait() && token.Error() != nil {
			glog.Warningf("device %s query device shadow data failed,request id = %s", device.base.Id, requestId)
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) UploadFile(filename string) AsyncResult {
	asyncResult := NewBooleanAsyncResult()
	go func() {
		// 构造获取文件上传URL的请求
		requestParas := iot.FileRequestServiceEventParas{
			FileName: filename,
		}

		serviceEvent := iot.FileRequestServiceEvent{
			Paras: requestParas,
		}
		serviceEvent.ServiceId = "$file_manager"
		serviceEvent.EventTime = iot.GetEventTimeStamp()
		serviceEvent.EventType = "get_upload_url"

		var services []iot.FileRequestServiceEvent
		services = append(services, serviceEvent)
		request := iot.FileRequest{
			Services: services,
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request)); token.Wait() && token.Error() != nil {
			glog.Warningf("publish file upload request url failed")
			asyncResult.completeError(&DeviceError{
				errorMsg: "publish file upload request url failed",
			})
			return
		}
		glog.Info("publish file upload request url success")

		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				_, ok := device.base.fileUrls[filename+topics.FileActionUpload]
				if ok {
					glog.Infof("platform send file upload url success")
					goto ENDFOR
				}

			}
		}
	ENDFOR:

		if len(device.base.fileUrls[filename+topics.FileActionUpload]) == 0 {
			glog.Errorf("get file upload url failed")
			asyncResult.completeError(&DeviceError{
				errorMsg: "get file upload url failed",
			})
			return
		}
		glog.Infof("file upload url is %s", device.base.fileUrls[filename+topics.FileActionUpload])

		//filename = smartFileName(filename)
		uploadFlag := CreateHttpClient().UploadFile(filename, device.base.fileUrls[filename+topics.FileActionUpload])
		if !uploadFlag {
			glog.Errorf("upload file failed")
			asyncResult.completeError(&DeviceError{
				errorMsg: "upload file failed",
			})
			return
		}

		response := CreateFileUploadDownLoadResultResponse(filename, topics.FileActionUpload, uploadFlag)

		token := device.base.Client.Publish(iot.FormatTopic(topics.PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(response))
		if token.Wait() && token.Error() != nil {
			glog.Error("report file upload file result failed")
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) DownloadFile(filename string) AsyncResult {
	asyncResult := NewBooleanAsyncResult()
	go func() {
		// 构造获取文件上传URL的请求
		requestParas := iot.FileRequestServiceEventParas{
			FileName: filename,
		}

		serviceEvent := iot.FileRequestServiceEvent{
			Paras: requestParas,
		}
		serviceEvent.ServiceId = "$file_manager"
		serviceEvent.EventTime = iot.GetEventTimeStamp()
		serviceEvent.EventType = "get_download_url"

		var services []iot.FileRequestServiceEvent
		services = append(services, serviceEvent)
		request := iot.FileRequest{
			Services: services,
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request)); token.Wait() && token.Error() != nil {
			glog.Warningf("publish file download request url failed")
			asyncResult.completeError(&DeviceError{
				errorMsg: "publish file download request url failed",
			})
			return
		}

		ticker := time.Tick(time.Second)
		for {
			select {
			case <-ticker:
				_, ok := device.base.fileUrls[filename+topics.FileActionDownload]
				if ok {
					glog.Infof("platform send file upload url success")
					goto ENDFOR
				}

			}
		}
	ENDFOR:

		if len(device.base.fileUrls[filename+topics.FileActionDownload]) == 0 {
			glog.Errorf("get file download url failed")
			asyncResult.completeError(&DeviceError{
				errorMsg: "get file download url failed",
			})
			return
		}

		downloadFlag := CreateHttpClient().DownloadFile(filename, device.base.fileUrls[filename+topics.FileActionDownload])
		if !downloadFlag {
			glog.Errorf("down load file { %s } failed", filename)
			asyncResult.completeError(&DeviceError{
				errorMsg: "down load file failedd",
			})
			return
		}

		response := CreateFileUploadDownLoadResultResponse(filename, topics.FileActionDownload, downloadFlag)

		token := device.base.Client.Publish(iot.FormatTopic(topics.PlatformEventToDeviceTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(response))
		if token.Wait() && token.Error() != nil {
			glog.Error("report file upload file result failed")
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}

	}()

	return asyncResult
}

func (device *asyncDevice) ReportDeviceInfo(swVersion, fwVersion string) AsyncResult {
	asyncResult := NewBooleanAsyncResult()
	go func() {
		event := iot.ReportDeviceInfoServiceEvent{
			BaseServiceEvent: iot.BaseServiceEvent{
				ServiceId: "$device_info",
				EventType: "device_info_report",
				EventTime: iot.GetEventTimeStamp(),
			},
			Paras: iot.ReportDeviceInfoEventParas{
				DeviceSdkVersion: iot.SdkInfo()["sdk-version"],
				SwVersion:        swVersion,
				FwVersion:        fwVersion,
			},
		}

		request := iot.ReportDeviceInfoRequest{
			DeviceId: device.base.Id,
			Services: []iot.ReportDeviceInfoServiceEvent{event},
		}

		token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request))
		if token.Wait() && token.Error() != nil {
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}

	}()

	return asyncResult
}

func (device *asyncDevice) ReportLogs(logs []iot.DeviceLog) AsyncResult {
	asyncresult := NewBooleanAsyncResult()

	go func() {
		var services []iot.ReportDeviceLogServiceEvent

		for _, logEntry := range logs {
			service := iot.ReportDeviceLogServiceEvent{
				BaseServiceEvent: iot.BaseServiceEvent{
					ServiceId: "$log",
					EventType: "log_report",
					EventTime: iot.GetEventTimeStamp(),
				},
				Paras: logEntry,
			}

			services = append(services, service)
		}

		request := iot.ReportDeviceLogRequest{
			Services: services,
		}

		fmt.Println(iot.Interface2JsonString(request))

		topic := iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id)

		token := device.base.Client.Publish(topic, 1, false, iot.Interface2JsonString(request))

		if token.Wait() && token.Error() != nil {
			glog.Errorf("device %s report log failed", device.base.Id)
			asyncresult.completeError(token.Error())
		} else {
			asyncresult.completeSuccess()
		}
	}()

	return asyncresult
}

func (device *asyncDevice) SetSubDevicesAddHandler(handler SubDevicesAddHandler) {
	device.base.subDevicesAddHandler = handler
}

func (device *asyncDevice) SetSubDevicesDeleteHandler(handler SubDevicesDeleteHandler) {
	device.base.subDevicesDeleteHandler = handler
}

func (device *asyncDevice) UpdateSubDeviceState(subDevicesStatus iot.SubDevicesStatus) AsyncResult {
	glog.Infof("begin to update sub-devices status")

	asyncResult := NewBooleanAsyncResult()

	go func() {
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
				EventTime: iot.GetEventTimeStamp(),
				Paras:     sds,
			}

			request := iot.EventData{
				DeviceId: device.base.Id,
				Services: []iot.ServiceEvent{requestEventService},
			}

			if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request)); token.Wait() && token.Error() != nil {
				glog.Warningf("gateway %s update sub devices status failed", device.base.Id)
				asyncResult.completeError(token.Error())
				return
			}
		}
		asyncResult.completeSuccess()
		glog.Info("gateway  update sub devices status failed", device.base.Id)
	}()

	return asyncResult
}

func (device *asyncDevice) DeleteSubDevices(deviceIds []string) AsyncResult {
	glog.Infof("begin to delete sub-devices %s", deviceIds)

	asyncResult := NewBooleanAsyncResult()

	go func() {
		subDevices := struct {
			Devices []string `json:"devices"`
		}{
			Devices: deviceIds,
		}

		requestEventService := iot.ServiceEvent{
			ServiceId: "$sub_device_manager",
			EventType: "delete_sub_device_request",
			EventTime: iot.GetEventTimeStamp(),
			Paras:     subDevices,
		}

		request := iot.EventData{
			DeviceId: device.base.Id,
			Services: []iot.ServiceEvent{requestEventService},
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request)); token.Wait() && token.Error() != nil {
			glog.Warningf("gateway %s delete sub devices request send failed", device.base.Id)
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}

		glog.Warningf("gateway %s delete sub devices request send success", device.base.Id)
	}()

	return asyncResult
}

func (device *asyncDevice) AddSubDevices(deviceInfos []iot.SubDeviceInfo) AsyncResult {
	asyncResult := NewBooleanAsyncResult()

	go func() {
		devices := struct {
			Devices []iot.SubDeviceInfo `json:"devices"`
		}{
			Devices: deviceInfos,
		}

		requestEventService := iot.ServiceEvent{
			ServiceId: "$sub_device_manager",
			EventType: "add_sub_device_request",
			EventTime: iot.GetEventTimeStamp(),
			Paras:     devices,
		}

		request := iot.EventData{
			DeviceId: device.base.Id,
			Services: []iot.ServiceEvent{requestEventService},
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(request)); token.Wait() && token.Error() != nil {
			glog.Warningf("gateway %s add sub devices request send failed", device.base.Id)
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}

		glog.Warningf("gateway %s add sub devices request send success", device.base.Id)
	}()

	return asyncResult
}

func (device *asyncDevice) SyncAllVersionSubDevices() AsyncResult {
	asyncResult := NewBooleanAsyncResult()

	go func() {
		dataEntry := iot.ServiceEvent{
			ServiceId: "$sub_device_manager",
			EventType: "sub_device_sync_request",
			EventTime: iot.GetEventTimeStamp(),
			Paras: struct {
			}{},
		}

		var dataEntries []iot.ServiceEvent
		dataEntries = append(dataEntries, dataEntry)

		data := iot.EventData{
			Services: dataEntries,
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(data)); token.Wait() && token.Error() != nil {
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}

func (device *asyncDevice) SyncSubDevices(version int) AsyncResult {
	asyncResult := NewBooleanAsyncResult()

	go func() {
		syncParas := struct {
			Version int `json:"version"`
		}{
			Version: version,
		}

		dataEntry := iot.ServiceEvent{
			ServiceId: "$sub_device_manager",
			EventType: "sub_device_sync_request",
			EventTime: iot.GetEventTimeStamp(),
			Paras:     syncParas,
		}

		var dataEntries []iot.ServiceEvent
		dataEntries = append(dataEntries, dataEntry)

		data := iot.EventData{
			Services: dataEntries,
		}

		if token := device.base.Client.Publish(iot.FormatTopic(topics.DeviceToPlatformTopic, device.base.Id), device.base.qos, false, iot.Interface2JsonString(data)); token.Wait() && token.Error() != nil {
			glog.Errorf("send sync sub device request failed")
			asyncResult.completeError(token.Error())
		} else {
			asyncResult.completeSuccess()
		}
	}()

	return asyncResult
}
