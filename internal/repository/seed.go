package repository

import (
	"database/sql"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
	logger "github.com/sirupsen/logrus"
)

func SeedData() (err error) {
	uri := config.ReadEnvString("DB_URI")

	db, _ := sql.Open(dbDriver, uri)
	defer db.Close()

	seedQueries := []string{
		//roles
		`INSERT INTO roles (id, name) VALUES (1, 'super admin')`,
		`INSERT INTO roles (id, name) VALUES (2, 'admin')`,
		`INSERT INTO roles (id, name) VALUES (3, 'user')`,

		//grades
		`INSERT INTO grades (id,name, points) VALUES (1, 'J1',1000)`,
		`INSERT INTO grades (id,name, points) VALUES (2,'J2',900)`,
		`INSERT INTO grades (id,name, points) VALUES (3,'J4',800)`,
		`INSERT INTO grades (id,name, points) VALUES (4,'J6',700)`,
		`INSERT INTO grades (id,name, points) VALUES (5,'J7',600)`,
		`INSERT INTO grades (id,name, points) VALUES (6,'J8',500)`,
		`INSERT INTO grades (id,name, points) VALUES (7,'J9',400)`,
		`INSERT INTO grades (id,name, points) VALUES (8,'J10',300)`,
		`INSERT INTO grades (id,name, points) VALUES (9,'J11',200)`,
		`INSERT INTO grades (id,name, points) VALUES (10,'J12',100)`,

		//corevalues
		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (1,'Trust','We foster trust by being transparent, reliable, and accountable in all our actions.',null)`,

		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (2,'Respect','We respect individual opinions and believe in living up to and exceeding our own standards.',null)`,

		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (3,'Integrity & Ethics','We uphold integrity in every action, ensuring our decisions align with the highest moral standards.',null)`,

		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (4,'Employee Focus','We create a sense of shared purpose and understand that the power of the collective and that of the team is bigger than the individual.',null)`,

		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (5, 'Customer Focus', 'We prioritize understanding and meeting our customers'' needs, exceeding expectations with every interaction.', null)`,

		`INSERT INTO core_values (id,name,description, parent_core_value_id) VALUES (6,'Technical Excellence','We are committed to delivering excellence in every product, service, and experience we provide, striving for continuous improvement.',null)`,

		//badges
		`INSERT INTO badges (id,name,reward_points) VALUES (1,'Bronze',1500)`,
		`INSERT INTO badges (id,name,reward_points) VALUES (2,'Silver',3000)`,
		`INSERT INTO badges (id,name,reward_points) VALUES (3,'Gold',5000)`,
		`INSERT INTO badges (id,name,reward_points) VALUES (4,'Platinum',7000)`,

		//users
		`INSERT INTO users (id,employee_id,first_name,last_name,email,password, designation,reward_quota_balance,role_id,grade_id)
		VALUES (1,'900','Sindhu','R', 'sindhu.r@joshsoftware.com','$2a$14$5smJBNxiWYKDy2WK0tzo7OggboToB/lr2jM9Q6qSI63gHMr08/212','Manager - HR',5000 , 2, 6)`,

		`INSERT INTO users (id,employee_id,first_name,last_name,email,password, designation,reward_quota_balance,role_id,grade_id)
		VALUES (2,'281','Saurabh','Gaji', 'saurabh.gaji@joshsoftware.com','$2a$14$5smJBNxiWYKDy2WK0tzo7OggboToB/lr2jM9Q6qSI63gHMr08/212','Senior Executive - HRBP',3000 , 2, 8)`,

		`INSERT INTO users (id,employee_id,first_name,last_name,email,password, designation,reward_quota_balance,role_id,grade_id)
		VALUES (3,'717','SriGayathriKavya','Ruttala','srigayathri.ruttala@joshsoftware.com','$2a$14$5smJBNxiWYKDy2WK0tzo7OggboToB/lr2jM9Q6qSI63gHMr08/212','Executive - HRBP',2000 , 2, 9)`,

		//organization config
		`INSERT INTO organization_config (id,reward_multiplier,reward_quota_renewal_frequency,timezone,created_by,updated_by) VALUES (1,10,1,'UTC',1,1)`,
	}

	for _, query := range seedQueries {
		_, err := db.Exec(query)
		if err != nil {
			logger.Error("Err", "failed to execute seed query (%s): %v", query, err)
			return apperrors.InternalServer
		}
	}

	logger.Info("Seed data loaded successfully.")
	return
}
