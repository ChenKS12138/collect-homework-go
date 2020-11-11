package auth

// Code auth code
type Code = uint32;


// Create W; Read R; Update R+W; Delete X+W;
const (
	// CodeFileX file excuse
	CodeFileX Code=0b1 << 0
	// CodeFileW file read
	CodeFileW Code=0b1 << 1
	// CodeFileR file write
	CodeFileR Code=0b1 << 2

	// CodeProjectX file excuse
	CodeProjectX Code=0b1 << 3
	// CodeProjectW project excuse
	CodeProjectW Code=0b1 << 4
	// CodeProjectR project write
	CodeProjectR Code=0b1 << 5
	
	// CodeAdminX admin excuse
	CodeAdminX Code=0b1 << 6
	// CodeAdminW admin read
	CodeAdminW Code=0b1 << 7
	// CodeAdminR admin write
	CodeAdminR Code=0b1 << 8
)

// VerifyAuthCode verify auth code
func VerifyAuthCode(src Code,target Code)bool{
	return src & target == target;
}