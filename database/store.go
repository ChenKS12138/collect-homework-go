package database

// Store store
var Store struct {
	Admin *AdminStore
	Project *ProjectStore
	Submission *SubmissionStore
	InvitationCode *InvitationCodeStore
}