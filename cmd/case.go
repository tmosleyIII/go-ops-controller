package cmd

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	caseInfo     = regexp.MustCompile(`^X([A-Z])([0-9]{2})([0-9]{1,7})`)
	origCaseType string
	caseYear     int
	seqNumber    int
)

type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

// Scan implements the Scanner interface.
func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface
func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

// Case represents the layout of cases in the system
type Case struct {
	ID                       int        `json:"id"`
	Barcode                  *string    `json:"barcode"`
	Biopsy_date              *time.Time `json:"biopsy_date"`
	Diagnosed_at             *time.Time `json:"diagnosed_at,omitempty"`
	Payment_method           *string    `json:"payment_method,omitempty"`
	Previous_case_number     *string    `json:"previous_case_number,omitempty"`
	Special_billing_request  *string    `json:"special_billing_request,omitempty"`
	Validated_at             *time.Time `json:"validated_at,omitempty"`
	Aetna_case               *bool      `json:"aetna_case"`
	Clearpath                *bool      `json:"clearpath"`
	Digital                  *bool      `json:"digital"`
	Digital_field_locked     *bool      `json:"digital_field_locked"`
	Orders_processing        *bool      `json:"orders_processing"`
	Rush                     *bool      `json:"rush"`
	Seq_number               *int64     `json:"seq_number"`
	Status                   *int64     `json:"status"`
	Year                     *int64     `json:"year"`
	Created_at               *time.Time `json:"created_at"`
	Updated_at               *time.Time `json:"updated_at"`
	Deleted_at               *time.Time `json:"deleted_at,omitempty"`
	Case_type_id             *string    `json:"case_type_id"`
	Contributor_id           *int64     `json:"contributor_id"`
	Lab_location_id          *int64     `json:"lab_location_id"`
	Location_id              *int64     `json:"location_id"`
	Patient_id               *int64     `json:"patient_id"`
	Previous_case_id         *int64     `json:"previous_case_id,omitempty"`
	Shipping_tracking_number *string    `json:"shipping_tracking_number,omitempty"`
	Order_number             *string    `json:"order_number,omitempty"`
	Consult_created_at       *time.Time `json:"consult_created_at,omitempty"`
	Slide_prep_ready         *bool      `json:"slide_prep_ready,omitempty"`
	Consult_pending          *bool      `json:"consult_pending"`
	Shipped_at               *time.Time `json:"shipped_at,omitempty"`
	Scope_status             *int64     `json:"scope_status,omitempty"`
	Second_step_complete     *bool      `json:"second_step_complete,omitempty"`
	Sent_iguana              *bool      `json:"sent_iguana,omitempty"`
	Scope_reviewed           *bool      `json:"scope_reviewed,omitempty"`
	Scope_tag                *string    `json:"scope_tag,omitempty"`
	Signed_out_at            *time.Time `json:"signed_out_at,omitempty"`
	Billing_completed_at     *time.Time `json:"billing_completed_at,omitempty"`
	Signed_out_by_id         *int64     `json:"signed_out_by_id,omitempty"`
	Previous_checked         *bool      `json:"previous_checked,omitempty"`
	Read_wsi                 *bool      `json:"read_wsi,omitempty"`
	Physician_assistant      *string    `json:"physician_assistant,omitempty"`
}

// handleError Receive error codes and display them
func handleError(w http.ResponseWriter, code int) {
	JSONError(w, Error{codes[code], code}, code)
}

// CaseShow Return case information for one case
func CaseShow(w http.ResponseWriter, r *http.Request) {
	prettyPrint := getPrettyPrintValue(r)
	accNumber := strings.ToUpper(getAccNumVar(r))
	if accNumber == "" {
		handleError(w, 400)
		return
	}

	origCaseType, caseYear, seqNumber = FindCaseSubstrings(accNumber)

	cases, err := GetCase(seqNumber, origCaseType, caseYear)
	switch {
	case err == sql.ErrNoRows:
		handleError(w, 404)
		return
	case err != nil:
		log.Fatal(err)
		handleError(w, 500)
		return
	}

	//tpl.ExecuteTemplate(w, "show_case.gohtml", cases)
	JSON(w, cases, prettyPrint, 200)
}

// CaseDelete Delete specified case
func CaseDelete(w http.ResponseWriter, r *http.Request) {
}

func FindCaseSubstrings(accNumber string) (string, int, int) {
	results := caseInfo.FindStringSubmatch(accNumber)
	cType := results[1]
	cYear := "20" + results[2]
	sNum := results[3]
	d, err := strconv.Atoi(cYear)
	if err != nil {
		log.Fatal(err)
	}

	s, err := strconv.Atoi(sNum)
	if err != nil {
		log.Fatal(err)
	}
	return cType, d, s
}

func GetCase(seq int, ctype string, cyear int) (*Case, error) {

	cases := &Case{}
	stmt, err := db.Prepare("SELECT * FROM public.case WHERE seq_number = $1 AND case_type_id = $2 AND year = $3")
	if err != nil {
		log.Fatal(err)
	}
	err = stmt.QueryRow(seq, ctype, cyear).Scan(&cases.ID, &cases.Barcode, &cases.Biopsy_date, &cases.Diagnosed_at, &cases.Payment_method, &cases.Previous_case_number, &cases.Special_billing_request, &cases.Validated_at, &cases.Aetna_case, &cases.Clearpath, &cases.Digital, &cases.Digital_field_locked, &cases.Orders_processing, &cases.Rush, &cases.Seq_number, &cases.Status, &cases.Year, &cases.Created_at, &cases.Updated_at, &cases.Deleted_at, &cases.Case_type_id, &cases.Contributor_id, &cases.Lab_location_id, &cases.Location_id, &cases.Patient_id, &cases.Previous_case_id, &cases.Shipping_tracking_number, &cases.Order_number, &cases.Consult_created_at, &cases.Slide_prep_ready, &cases.Consult_pending, &cases.Shipped_at, &cases.Scope_status, &cases.Second_step_complete, &cases.Sent_iguana, &cases.Scope_reviewed, &cases.Scope_tag, &cases.Signed_out_at, &cases.Billing_completed_at, &cases.Signed_out_by_id, &cases.Previous_checked, &cases.Read_wsi, &cases.Physician_assistant)
	if err == sql.ErrNoRows {
		return nil, err
	}

	return cases, err
}
