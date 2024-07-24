package script

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joshsoftware/peerly-backend/internal/pkg/apperrors"
	"github.com/joshsoftware/peerly-backend/internal/pkg/config"
)

type GetUserListRespData struct {
	Data []IntranetUserData `json:"data"`
}
type IntranetUserData struct {
	Id             int            `json:"id"`
	Email          string         `json:"email"`
	PublicProfile  PublicProfile  `json:"public_profile"`
	EmpolyeeDetail EmployeeDetail `json:"employee_detail"`
}
type PublicProfile struct {
	ProfileImgUrl string `json:"profile_image_url"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
}
type EmployeeDetail struct {
	EmployeeId  string      `json:"employee_id"`
	Designation Designation `json:"designation"`
	Grade       string      `json:"grade"`
}
type Designation struct {
	Name string `json:"name"`
}

const (
	POST         = "POST"
	GET          = "GET"
	IntranetAuth = "Intranet-Auth"
)

func LoadUserScript() error {

	// fmt.Printf("in LoadUserScript 45")

	scriptErr := apperrors.InternalServerError

	client := &http.Client{}
	url := fmt.Sprintf(config.PeerlyBaseUrl()+"/users?page=%d", 1)
	req, err := http.NewRequest(GET, url, nil)
	if err != nil {
		fmt.Printf("Error in creating new request, err %+v", err)
		return scriptErr
	}

	req.Header.Add(IntranetAuth, config.IntranetAuthToken())
	resp, err := client.Do(req)
	// fmt.Println("first request done")
	if err != nil {
		fmt.Println("Error getuserlist from intranet api: ", err.Error())
		return scriptErr
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Intranet get api failed with status code: ", resp.StatusCode)
		return scriptErr
	}
	defer resp.Body.Close()

	var respData GetUserListRespData

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error in io.readall. err: ", err.Error())
		return scriptErr
	}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Println("Error in unmarshalling data, err: ", err.Error())
		return scriptErr
	}

	data := respData.Data

	for i := 0; i < len(data); i++ {
		client := &http.Client{}
		url := fmt.Sprintf(config.PeerlyBaseUrl() + "/user/register")
		jsonData, err := json.Marshal(data[i])
		if err != nil {
			fmt.Println("err in json data marshalling")
			return scriptErr
		}
		req, err := http.NewRequest(POST, url, bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error in creating new request, err %+v", err)
			return scriptErr
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error in resp", err.Error())
		}
		if resp.StatusCode != http.StatusOK {
			switch resp.StatusCode {
			case http.StatusConflict:
				fmt.Println("User already exists!")
			case http.StatusBadRequest:
				fmt.Println("Incomplete user details")
			default:
				fmt.Println("Error statuscode: ", resp.StatusCode)
			}

		}
	}

	return nil
}
