package types

type Employee struct {
	Name           string
	Interests      []string
	AssignedToGift string
}

type Gift struct {
	Name               string
	Categories         []string
	AssignedToEmployee string
}

type Config struct {
	EmployeesPath string
	GiftsPath     string
}

type AssignEmployeeToGiftRequest struct {
	EmployeeName string `json:"employee_name"`
}

type AssignEmployeeToGiftResponse struct {
	Error string `json:"error"`
}

type Loader interface {
	LoadEmployees(es []Employee, err error)
	LoadGifts(es []Gift, err error)
}

type Storage interface {
	AssignEmployeeToGift(employeeName, giftName string) (err error)
	GetEmployees() (es []Employee, err error)
	GetUnassignedEmployees() (es []Employee, err error)
	GetEmployeeByName(name string) (e Employee, err error)
	GetGifts() (gs []Gift, err error)
	GetUnassignedGifts() (gs []Gift, err error)
	GetGiftByName(name string) (g Gift, err error)
}

type Server interface {
	Serve() (err error)
	Close() (err error)
}

type Dependencies struct {
	Storage Storage
	Server Server
}
