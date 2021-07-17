package auth

// Code auth code
type Code = uint32

// Create W; Read R; Update R+W; Delete X+W;
const (
	// CodeFileX file excuse
	CodeFileX Code = 0b1 << iota
	// CodeFileW file read
	CodeFileW
	// CodeFileR file write
	CodeFileR

	// CodeProjectX file excuse
	CodeProjectX
	// CodeProjectW project excuse
	CodeProjectW
	// CodeProjectR project write
	CodeProjectR

	// CodeAdminX admin excuse
	CodeAdminX
	// CodeAdminW admin read
	CodeAdminW
	// CodeAdminR admin write
	CodeAdminR
)

// VerifyAuthCode verify auth code
func VerifyAuthCode(src Code, target Code) bool {
	return src&target == target
}
