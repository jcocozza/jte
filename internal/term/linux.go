//go:build linux || freebsd || netbsd || openbsd

package term

import (
	"syscall"
	"unsafe"
)

const (
	ECHO   = 0x00000008
	TCGETS = 0x5401
	TCSETS = 0x5402
)

const (
	ioctlReadTermios  = syscall.TCGETS
	ioctlWriteTermios = syscall.TCSETS
)


func enableRawMode() (*RawMode, error) {
	// Open /dev/tty explicitly instead of using stdin
	//fd, err := syscall.Open("/dev/tty", syscall.O_RDWR, 0)
	//if err != nil {
	//	return nil, err
	//}
	fd := syscall.Stdin

	var termios syscall.Termios
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(ioctlReadTermios),
		uintptr(unsafe.Pointer(&termios)),
	)
	if errno != 0 {
		syscall.Close(fd)
		return nil, errno
	}

	// Save original state
	originalTermios := termios

	// Modify the settings:
	// - disable ECHO
	// - disable CANON (canonical mode)
	// - disable kill signals (ctrl-c and ctrl-z)
	// - disable ctrl-v
	termios.Lflag &^= syscall.ECHO | syscall.ICANON | syscall.ISIG | syscall.IXON | syscall.IEXTEN | syscall.ICRNL
	// - disable software flow control (ctrl-s and ctrl-q)
	// - fix ctrl-m
	// - some other, not important flags
	termios.Iflag &^= syscall.IXON | syscall.ICRNL | syscall.BRKINT | syscall.INPCK | syscall.ISTRIP
	// turn off output processing
	termios.Oflag &^= syscall.OPOST
	// disable other flags
	termios.Cflag &^= syscall.CS8

	// set a timeout for reading
	//termios.Cc[syscall.VMIN] = 0
	//termios.Cc[syscall.VTIME] = 10

	_, _, errno = syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(fd),
		uintptr(ioctlWriteTermios),
		uintptr(unsafe.Pointer(&termios)),
	)
	if errno != 0 {
		syscall.Close(fd)
		return nil, errno
	}

	return &RawMode{
		originalState: originalTermios,
		fd:            fd,
	}, nil
}

func restore(r *RawMode) error {
	if r == nil {
		return nil
	}
	termios, ok := r.originalState.(syscall.Termios)
	if !ok {
		return nil
	}
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		uintptr(r.fd),
		uintptr(ioctlWriteTermios),
		uintptr(unsafe.Pointer(&termios)),
	)
	// Always try to close the fd
	syscall.Close(r.fd)
	if errno != 0 {
		return errno
	}
	return nil
}

// todo: this is not guarenteed to work; see https://viewsourcecode.org/snaptoken/kilo/03.rawInputAndOutput.html#window-size-the-hard-way
// returns row, col, err
func getWindowSize() (int, int, error) {
	var ws struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdout),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(&ws)))
	if errno != 0 {
		return -1, -1, errno
	}
	return int(ws.Row), int(ws.Col), nil
}
