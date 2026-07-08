//
// 
package ports

//
type Orchestrator interface {
    RunCLI(url string)
    RunHTTP(url string)
}
