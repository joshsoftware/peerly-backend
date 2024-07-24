package constants

var PerPage = 400

type UserIdCtxKey string
type RoleCtxKey string

var UserId UserIdCtxKey = "userId"
var Role RoleCtxKey = "role"

const AuthorizationHeader string = "Authorization"
const ClientCode = "Client-Code"
const UserRole = "user"
const DefaultPageSize = 400
const UserId = "userId"
const Role = "role"
const IntranetAuth = "Intranet-Auth"
const PeerlyValidationPath = "/api/peerly/v1/sessions/login"
const GetIntranetUserDataPath = "/api/peerly/v1/users/"
const ListIntranetUsersPath = "/api/peerly/v1/users?page=%d&per_page=%d"
