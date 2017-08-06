# Notes
Something wrong with the marshal / unmarshal logic
client.ReadDir fails because of server?

- Connect and download license => works, download again => does not work
- attrs test does not test for dir?!

# Must-Changes
func NewRequest(method, path string) Request {
	// XXX
	request := Request{Method: method, Filepath: filepath.ToSlash(filepath.Clean(path))}



# Help
search for " filepath.Clean", replace with " filepath.ToSlash(filepath.Clean"

# Interesting parts

func fileinfo(h FileInfoer, r Request) (responsePacket, error) {

}
