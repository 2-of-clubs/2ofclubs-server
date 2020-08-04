package handler

import (
	"fmt"
	"github.com/2-of-clubs/2ofclubs-server/app/model"
	"github.com/go-playground/validator"
	"gorm.io/gorm"
	"net/http"
	"strings"
)


func GetClubs(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	status := model.NewStatus()

	type Club struct {
		club_id int
	}
	clubs := []Club{}
	status.Message = model.ClubsFound
	activeTags := flatten(filterTags(extractTags(db, r)))
	//for _, v := range activeTags {
	//fmt.Println(v.Name)
	//}
	fmt.Println(activeTags)
	db.Table(model.ClubTagTable).Where("tag_name IN ?", activeTags).Find(&clubs)
	fmt.Println(clubs)
	WriteData(GetJSON(status), http.StatusOK, w)
}

/*
Check if the email & username is available (RecordExists)
*/
func CreateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	/*
		-> Validate User token (Done)
		-> Get Token Claims (User must exist then, unless deleted?) (You can put a check on Record Exists on the deleted column as long as it's null it'll exist then) (Done)
		-> Extract Username and return User struct (Done)
		-> Check if the club is available (Done)
		-> Extract New Club Info to Struct (Done)
		-> Insert to User Manages (Done)
		-> Update user (Done)
		-> Set the user as the owner of the club (Done)
	*/
	claims := GetTokenClaims(r)
	user := model.NewUser()
	uname := fmt.Sprintf("%v", claims["sub"])
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, uname, user)
	club := model.NewClub()
	extractBody(r, club)
	validate := validator.New()
	err := validate.Struct(club)
	clubExists := SingleRecordExists(db, model.ClubTable, model.NameColumn, club.Name, model.NewClub())
	emailExists := SingleRecordExists(db, model.ClubTable, model.EmailColumn, club.Email, model.NewClub())
	status := model.NewStatus()
	// Keeping userExists as a check even though the user should exist given the valid token because there's a chance that the user is deleted
	// In this case the user will still exist in the database but will be inaccessible.
	if !emailExists && !clubExists && userExists && err == nil {
		db.Model(user).Association(model.ManagesColumn).Append(club)
		db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ? AND is_owner = ?", user.ID, club.ID, false).Update(model.IsOwnerColumn, true)
		status.Message = model.SuccessClubCreation
		status.Code = model.SuccessCode
	} else {
		status.Message = model.FailureClubCreation
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}

//func getClubInfo(r *http.Request) *model.Club {
//	decoder := json.NewDecoder(r.Body)
//	club := model.NewClub()
//	decoder.Decode(club)
//	return club
//}

func GetClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getClubInfo(db, w, r, model.AllClubInfo)
}
func getClubInfo(db *gorm.DB, w http.ResponseWriter, r *http.Request, infoType string) {
	var statusCode int
	var data string
	clubID := getVar(r, model.ClubIDVar)
	status := model.NewStatus()
	club := model.NewClub()
	found := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if !found {
		status.Message = model.ClubNotFound
	} else {
		switch strings.ToLower(infoType) {
		case model.AllClubInfo:
			clubDisplay := club.Display()
			loadClubData(db, club, clubDisplay)
			status.Data = clubDisplay
		case model.AllClubEventsHost:
			clubEvents := make(map[string][]model.Event)
			db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
			clubEvents[model.HostsColumn] = club.Hosts
			status.Data = clubEvents
		}
		status.Message = model.ClubFound
		status.Code = model.SuccessCode
	}
	statusCode = http.StatusOK
	data = GetJSON(status)
	WriteData(data, statusCode, w)
}

func GetClubEvents(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	getClubInfo(db, w, r, model.AllClubEventsHost)
}
func loadClubData(db *gorm.DB, club *model.Club, clubDisplay *model.ClubDisplay) {
	db.Table(model.ClubTable).Preload(model.SetsColumn).Find(club)
	db.Table(model.ClubTable).Preload(model.HostsColumn).Find(club)
	clubDisplay.Tags = flatten(filterTags(club.Sets))
	clubDisplay.Hosts = club.Hosts
}


func UpdateClub(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update Club")
}

/*
Returns true if the user is an owner of the club, false otherwise
*/
func isOwner(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return userClub.IsOwner
}

func isManager(db *gorm.DB, user *model.User, club *model.Club) bool {
	userClub := model.NewUserClub()
	res := db.Table(model.UserClubTable).Where("user_id = ? AND club_id = ?", user.ID, club.ID).First(userClub)
	return res.Error == nil
}

func RemoveManager(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	editManagers(db, w, r, model.OpRemove)
}

func AddManager(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	editManagers(db, w, r, model.OpAdd)
}

func UpdateClubTags(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var httpStatus int
	status := model.NewStatus()
	clubID := getVar(r, model.ClubIDVar)
	club := model.NewClub()
	claims := GetTokenClaims(r)
	username := fmt.Sprintf("%v", claims["sub"])
	user := model.NewUser()
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	userExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, username, user)
	// Must check with both user and club existing in the event that a user gets deleted but you manage to get a hold of their access token
	if userExists && clubExists && isManager(db, user, club) {
		tags := filterTags(extractTags(db, r))
		db.Model(club).Association(model.SetsColumn).Replace(tags)
		status.Message = model.TagsUpdated
		httpStatus = http.StatusOK
		status.Code = model.SuccessCode
	} else {
		status.Message = http.StatusText(http.StatusForbidden)
		httpStatus = http.StatusForbidden
	}
	WriteData(GetJSON(status), httpStatus, w)
}

/*
Adding or removing managers and their associations to a particular club
*/
func editManagers(db *gorm.DB, w http.ResponseWriter, r *http.Request, op string) {
	// Default messages set to manager addition, otherwise manager removal
	var successMessage = model.SuccessManagerAddition
	var failureMessage = model.FailureManagerAddition
	if op == model.OpRemove {
		successMessage = model.SuccessManagerRemove
		failureMessage = model.FailureManagerRemove
	}
	status := model.NewStatus()
	claims := GetTokenClaims(r)
	clubOwnerUsername := fmt.Sprintf("%v", claims["sub"])
	newManagerUname := getVar(r, model.UsernameVar)
	clubID := getVar(r, model.ClubIDVar)
	owner := model.NewUser()
	newManager := model.NewUser()
	club := model.NewClub()
	// Added user must exist
	ownerExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, clubOwnerUsername, owner)
	// If owner is found, then the owner struct isn't populated, which gives ID=0, but ID's start at 1, so this shouldn't cause any potential security issues
	managerExists := SingleRecordExists(db, model.UserTable, model.UsernameColumn, newManagerUname, newManager)
	clubExists := SingleRecordExists(db, model.ClubTable, model.IDColumn, clubID, club)
	if ownerExists && managerExists && clubExists {
		if isOwner(db, owner, club) && owner.Username != newManager.Username {
			switch op {
			case model.OpAdd:
				db.Model(newManager).Association(model.ManagesColumn).Append(club)
			case model.OpRemove:
				db.Model(newManager).Association(model.ManagesColumn).Delete(club)
			}
			status.Message = successMessage
			status.Code = model.SuccessCode
		} else {
			status.Message = failureMessage
		}
	} else {
		status.Message = failureMessage
	}
	WriteData(GetJSON(status), http.StatusOK, w)
}
