package docker

type ContainerInfo struct {
	ID    string
	Names []string
	Image string

	CPUPercent float64
	Mem        uint64
}
