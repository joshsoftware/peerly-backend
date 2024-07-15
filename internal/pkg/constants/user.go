package constants

var AuthorizationHeader string = "Authorization"
var ClientCode = "Client-Code"
var UserRole = "user"
var PerPage = 400

type UserIdCtxKey string
type RoleCtxKey string

var UserId UserIdCtxKey = "userId"
var Role RoleCtxKey = "role"

var IntranetAuth = "Intranet-Auth"
