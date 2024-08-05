package constants


var AppreciationColumns = []string{
	"id","core_value_id", "description","is_valid","total_rewards","quarter","sender","receiver","created_at","updated_at",
}
var CreateAppreciationColumns = []string{
	"core_value_id", "description","quarter","sender","receiver",
}

var OrgConfigColumns = []string{
	"id",
	"reward_multiplier",
	"reward_quota_renewal_frequency",
	"timezone",
	"created_by",
	"updated_by",
}