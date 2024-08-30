package constants

type UserIdCtxKey string
type RoleCtxKey string

// System Constants used to setup environment and basic functionality
const (
	AppName                = "APP_NAME"
	AppPort                = "APP_PORT"
	JWTSecret              = "JWT_SECRET"
	JWTExpiryDurationHours = "JWT_EXPIRY_DURATION_HOURS"
	DBURI                  = "DB_URI"
	IntranetClientCode     = "INTRANET_CLIENT_CODE"
	MigrationFolderPath    = "MIGRATION_FOLDER_PATH"
	IntranetAuthToken      = "INTRANET_AUTH_TOKEN"
	PeerlyBaseUrl          = "PEERLY_BASE_URL"
	IntranetBaseUrl        = "INTRANET_BASE_URL"
	POST                   = "POST"
	GET                    = "GET"
)

const (
	Admin int = iota
	User
)

// User required constants
const (
	AuthorizationHeader                  = "Authorization"
	ClientCode                           = "Client-Code"
	UserRole                             = "user"
	AdminRole                            = "admin"
	UserId                  UserIdCtxKey = "userId"
	Role                    RoleCtxKey   = "role"
	IntranetAuth                         = "Intranet-Auth"
	PeerlyValidationPath                 = "/api/peerly/v1/sessions/login"
	GetIntranetUserDataPath              = "/api/peerly/v1/users/"
	ListIntranetUsersPath                = "/api/peerly/v1/users?page=%d&per_page=%d"
)

// Pagination Metadata constants
const (
	DefaultPageNumber = 1
	DefaultPageSize   = 400
	MaxPageSize       = 1000
)

// Table Names
const (
	AppreciationsTable      = "appreciations"
	RewardsTable            = "rewards"
	UsersTable              = "users"
	CoreValuesTable         = "core_values"
	GradesTable             = "grades"
	OrganizationConfigTable = "organization_config"
	BadgeTable              = "badges"
	RolesTable              = "roles"
)

const DefaultOrgID = 1

const CheckIconLogo = "/peerly/assets/checkIcon.png"
