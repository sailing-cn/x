package jwt

const (
	Tenant     = "tenant"      //租户
	TenantId   = "tenant_id"   //租户标识
	TenantCode = "tenant_code" //租户编码
	EmployeeId = "employee_id" //员工标识
	AccountId  = "sub"         //账户标识
	Name       = "name"        //
	AppCode    = "app_code"    // 应用编码
	AppName    = "app_name"    // 应用名称
	AppId      = "app_id"
	FamilyName = "family_name" //
)

var (
	CLAIMS []string = []string{Tenant, TenantId, TenantCode,
		EmployeeId, AccountId, Name, AppCode, AppName, AppId, FamilyName}
)
