package jwt

const (
	Tenant     = "tenant"      //租户
	TenantId   = "tenant_id"   //租户标识
	TenantCode = "tenant_code" //租户编码
	EmployeeId = "employee_id" //员工标识
	AccountId  = "sub"         //账户标识
	Name       = "name"        //
)

var (
	CLAIMS []string = []string{Tenant, TenantId, TenantCode, EmployeeId, AccountId, Name}
)
