package sftp

import (
	"bytes"
	"os"
	"reflect"
	"testing"
	"time"
)

// ensure that attrs implemenst os.FileInfo
var _ os.FileInfo = new(fileInfo)

var unmarshalAttrsTests = []struct {
	b    []byte
	want *fileInfo
	rest []byte
}{
	{marshal(nil, struct{ Flags uint32 }{}), &fileInfo{mtime: time.Unix(int64(0), 0)}, nil},
	{marshal(nil, struct {
		Flags uint32
		Size  uint64
	}{ssh_FILEXFER_ATTR_SIZE, 20}), &fileInfo{size: 20, mtime: time.Unix(int64(0), 0)}, nil},
	{marshal(nil, struct {
		Flags       uint32
		Size        uint64
		Permissions uint32
	}{ssh_FILEXFER_ATTR_SIZE | ssh_FILEXFER_ATTR_PERMISSIONS, 20, 0644}), &fileInfo{size: 20, mode: os.FileMode(0644), mtime: time.Unix(int64(0), 0)}, nil},
	{marshal(nil, struct {
		Flags                 uint32
		Size                  uint64
		UID, GID, Permissions uint32
	}{ssh_FILEXFER_ATTR_SIZE | ssh_FILEXFER_ATTR_UIDGID | ssh_FILEXFER_ATTR_UIDGID | ssh_FILEXFER_ATTR_PERMISSIONS, 20, 1000, 1000, 0644}), &fileInfo{size: 20, mode: os.FileMode(0644), mtime: time.Unix(int64(0), 0)}, nil},
}

func TestUnmarshalAttrs(t *testing.T) {
	for _, tt := range unmarshalAttrsTests {
		stat, rest := unmarshalAttrs(tt.b)
		got := fileInfoFromStat(stat, "")
		tt.want.sys = got.Sys()
		if !reflect.DeepEqual(got, tt.want) || !bytes.Equal(tt.rest, rest) {
			t.Errorf("unmarshalAttrs(%#v): want %#v, %#v, got: %#v, %#v", tt.b, tt.want, tt.rest, got, rest)
		}
	}
}


//func TestMarshalAttrsForDirectory() {
//	fi := &fileInfo{name: "examples", size: 306, mode:2147484141, }
//
//	flags := 	var flags uint32 = ssh_FILEXFER_ATTR_SIZE |
//		ssh_FILEXFER_ATTR_PERMISSIONS |
//		ssh_FILEXFER_ATTR_ACMODTIME
//
//	got := marshal(flags, fi)
//	want := []byte{}
//	if !bytes.Equal(want, )
//	// win  => fi: &{name:examples sys:{FileAttributes:16 CreationTime:{LowDateTime:600060206 HighDateTime:30606641} LastAccessTime:{LowDateTime:1533014309 HighDateTime:30606642} LastWriteTime:{LowDateTime:1533014309 HighDateTime:30606642} FileSizeHigh:0 FileSizeLow:4096} pipe:false Mutex:{state:0 sema:0} path:C:\Users\mediacenter\go\src\github.com\sandreas\sftp\examples vol:0 idxhi:0 idxlo:0}
//
//	// unix => fi: &{name:examples size:306 mode:2147484141 modTime:{sec:63637652023 nsec:0 loc:0x1312560} sys:{Dev:16777220 Mode:16877 Nlink:9 Ino:14642480 Uid:501 Gid:20 Rdev:0 Pad_cgo_0:[0 0 0 0] Atimespec:{Sec:1502055237 Nsec:0} Mtimespec:{Sec:1502055223 Nsec:0} Ctimespec:{Sec:1502055223 Nsec:0} Birthtimespec:{Sec:1502055144 Nsec:0} Size:306 Blocks:0 Blksize:4096 Flags:0 Gen:0 Lspare:0 Qspare:[0 0]}}
//
//}

func TestFileStatFromInfoForDir(t *testing.T) {
	fi, _ := os.Stat("examples")
	flags, fileStat := fileStatFromInfo(fi)

	var got []byte
	got = marshalUint32(got, flags)
	if flags&ssh_FILEXFER_ATTR_SIZE != 0 {
		got = marshalUint64(got, fileStat.Size)
	}
	if flags&ssh_FILEXFER_ATTR_UIDGID != 0 {
		got = marshalUint32(got, fileStat.UID)
		got = marshalUint32(got, fileStat.GID)
	}
	if flags&ssh_FILEXFER_ATTR_PERMISSIONS != 0 {
		got = marshalUint32(got, fileStat.Mode)
	}
	if flags&ssh_FILEXFER_ATTR_ACMODTIME != 0 {
		got = marshalUint32(got, fileStat.Atime)
		got = marshalUint32(got, fileStat.Mtime)
	}

	want := []byte{
		0, 0, 0, 8, 101, 120, 97, 109, 112, 108, 101, 115, 0, 0, 0, 64, 100, 114, 119, 120, 114, 45, 120, 114, 45, 120, 32, 32, 32, 32, 57, 32, 53, 48, 49, 32, 32, 32, 32, 32, 32, 50, 48, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 32, 51, 48, 54, 32, 65, 117, 103, 32, 32, 54, 32, 50, 51, 58, 51, 51, 32, 101, 120, 97, 109, 112, 108, 101, 115, 0, 0, 0, 15, 0, 0, 0, 0, 0, 0, 1, 50, 0, 0, 1, 245, 0, 0, 0, 20, 0, 0, 65, 237, 89, 135, 139, 55, 89, 135, 139, 55,
	}

	if  !bytes.Equal(want, got) {
		t.Errorf("\nwant: %#v\n got: %#v", want, got)
	}
}

