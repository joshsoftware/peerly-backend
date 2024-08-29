package constants


type UserIdCtxKey string
type RoleCtxKey string
type RequestIDCtxKey string

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
	DeveloperKey           = "DEVELOPER_KEY"
)

// User required constants
const (
	RequestID               RequestIDCtxKey = "RequestID"
	AuthorizationHeader                     = "Authorization"
	ClientCode                              = "Client-Code"
	UserRole                                = "user"
	AdminRole                               = "admin"
	UserId                  UserIdCtxKey    = "userId"
	Role                    RoleCtxKey      = "role"
	IntranetAuth                            = "Intranet-Auth"
	PeerlyValidationPath                    = "/api/peerly/v1/sessions/login"
	GetIntranetUserDataPath                 = "/api/peerly/v1/users/"
	ListIntranetUsersPath                   = "/api/peerly/v1/users?page=%d&per_page=%d"
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

// EmailTemplate Icon url
const (
	BronzeBadgeIconImagePath        = "/peerly/assets/bronzeBadge.png"
	SilverBadgeIconImagePath        = "/peerly/assets/silverBadge.png"
	GoldBadgeIconImagePath          = "/peerly/assets/goldBadge.png"
	PlatinumIconImagePath           = "/peerly/assets/platinumBadge.png"
	CheckIconImagePath              = "/peerly/assets/checkIcon.png"
	ClosedEnvelopeIconImagePath     = "/peerly/assets/closedEnvelopeIcon.png"
	OpenEnvelopeIconImagePath       = "/peerly/assets/openEnvelopeIcon.png"
	RewardQuotaRenewalIconImagePath = "/peerly/assets/rewardQuotaRenewal.png"
)

//notificatio service account key file
const ServiceAccountKey = "serviceAccountKey.json"

// Email Dl group of HRs
const HRDLGroup = "dl_peerly.support@joshsoftware.com"