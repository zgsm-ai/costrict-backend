package chserver

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

/**
 * PortAllocator manages available ports allocation
 * @description
 * - Maintains port pool with unused/allocated/occupied states
 * - Provides thread-safe port allocation
 * - Supports port lookup by clientId and appName
 */
type PortAllocator struct {
	mu      sync.Mutex
	names   map[string]*PortAllocation // key: "clientId-appName" -> alloc
	ports   map[int]*PortAllocation    // port number -> alloc
	minPort int
	maxPort int
}

/**
 * NewPortAllocator creates a new port allocator instance
 * @param {int} minPort - Start port of allocation range
 * @param {int} maxPort - End port of allocation range
 * @returns {*PortAllocator} Initialized port allocator
 */
func NewPortAllocator(minPort, maxPort int) *PortAllocator {
	return &PortAllocator{
		names:   make(map[string]*PortAllocation),
		ports:   make(map[int]*PortAllocation),
		minPort: minPort,
		maxPort: maxPort,
	}
}

/**
 * AllocatePort assigns a new port mapping
 * @param {string} clientId - Client identifier
 * @param {string} appName - Application name
 * @param {int} clientPort - Client port number
 * @returns {PortAllocation, error} Allocation details or error
 */
func (pa *PortAllocator) AllocatePort(clientId, userId, appName string, clientPort int) (PortAllocation, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	now := time.Now().Local()
	key := clientId + "-" + userId + "-" + appName
	if alloc, exists := pa.names[key]; exists {
		alloc.ClientPort = clientPort
		alloc.Status = Allocated
		alloc.AllocTime = &now
		alloc.StartTime = nil
		return *alloc, nil
	}
	//先分配空槽
	for port := pa.minPort; port <= pa.maxPort; port++ {
		_, exists := pa.ports[port]
		if !exists {
			alloc := &PortAllocation{
				ClientId:    clientId,
				UserId:      userId,
				AppName:     appName,
				ClientPort:  clientPort,
				MappingPort: port,
				AllocTime:   &now,
				StartTime:   nil,
				Status:      Allocated,
			}
			pa.names[key] = alloc
			pa.ports[port] = alloc
			return *alloc, nil
		}
	}
	// 再分配回收再利用的旧槽
	for port := pa.minPort; port <= pa.maxPort; port++ {
		alloc, exists := pa.ports[port]
		if exists && alloc.Status == Freed {
			alloc = &PortAllocation{
				ClientId:    clientId,
				UserId:      userId,
				AppName:     appName,
				ClientPort:  clientPort,
				MappingPort: port,
				AllocTime:   &now,
				StartTime:   nil,
				Status:      Allocated,
			}
			pa.names[key] = alloc
			pa.ports[port] = alloc
			return *alloc, nil
		}
	}

	return PortAllocation{}, errors.New("no available ports")
}

func (pa *PortAllocator) applyNew(clientId, userId, appName string, clientPort, mappingPort int) (*PortAllocation, error) {
	//	先找该客户端的分配记录
	now := time.Now().Local()
	key := clientId + "-" + userId + "-" + appName
	if alloc, exists := pa.names[key]; exists {
		if clientPort != alloc.ClientPort {
			return nil, fmt.Errorf("client port conflict: %d - %d", clientPort, alloc.ClientPort)
		}
		if mappingPort != alloc.MappingPort {
			return nil, fmt.Errorf("mapping port conflict: %d - %d", mappingPort, alloc.MappingPort)
		}
		alloc.Status = Connected
		alloc.StartTime = &now
		return alloc, nil
	}
	//	如果没找到，说明之前没申请过，可能是因为cotund重启，导致客户端使用原端口重新连接
	//	此时只需要保证该端口还未被占用，即可直接给该客户端使用了
	if p, exists := pa.ports[mappingPort]; exists {
		return nil, fmt.Errorf("mapping port already allocated to: %+v", p)
	}

	alloc := &PortAllocation{
		ClientId:    clientId,
		UserId:      userId,
		AppName:     appName,
		ClientPort:  clientPort,
		MappingPort: mappingPort,
		AllocTime:   &now,
		StartTime:   &now,
		Status:      Connected,
	}
	pa.names[key] = alloc
	pa.ports[mappingPort] = alloc
	return alloc, nil
}

func (pa *PortAllocator) applyOld(clientPort, mappingPort int) (*PortAllocation, error) {
	alloc, exists := pa.ports[mappingPort]
	if !exists {
		return nil, fmt.Errorf("port [%d->%d] not exist", clientPort, mappingPort)
	}
	if clientPort != alloc.ClientPort {
		return nil, fmt.Errorf("port [%d] 'clientPort' conflict: [%d ~ %d]", alloc.MappingPort, clientPort, alloc.ClientPort)
	}
	if alloc.Status != Allocated {
		return nil, fmt.Errorf("port [%d] already used: %+v", alloc.MappingPort, alloc)
	}
	now := time.Now().Local()
	alloc.StartTime = &now
	alloc.Status = Connected
	return alloc, nil
}

/**
 *	隧道连接成功，应用已经分配的端口
 *	旧版本的cotun客户端，只有两个有效参数clientPort, mappingPort
 *	新版本的cotun客户端，会在http请求头中携带X-Client-Id, X-User-Id, X-App-Name
 */
func (pa *PortAllocator) OnConnected(c *PortAllocation, clientPort, mappingPort int) (*PortAllocation, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	if c.ClientId != "" && c.UserId != "" && c.AppName != "" {
		return pa.applyNew(c.ClientId, c.UserId, c.AppName, clientPort, mappingPort)
	}
	return pa.applyOld(clientPort, mappingPort)
}

/**
 *	隧道连接断开，端口仍然给该客户端保留10min
 */
func (pa *PortAllocator) OnDisconnected(alloc *PortAllocation) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	now := time.Now().Local()
	alloc.Status = Allocated
	alloc.AllocTime = &now
	alloc.StartTime = nil
}

/**
 * FreePort frees allocated port
 * @param {string} clientId - Client identifier
 * @param {string} appName - Application name (empty releases all client ports)
 */
func (pa *PortAllocator) FreePort(clientId, userId, appName string) []PortAllocation {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	allocs := []PortAllocation{}
	if appName == "" || userId == "" {
		for _, alloc := range pa.names {
			if alloc.Status == Freed {
				continue
			}
			if alloc.ClientId == clientId {
				alloc.Status = Freed
				allocs = append(allocs, *alloc)
			}
		}
		return allocs
	}

	key := clientId + "-" + userId + "-" + appName
	if alloc, exists := pa.names[key]; exists {
		alloc.Status = Freed
		allocs = append(allocs, *alloc)
	}
	for _, a := range allocs {
		key := a.ClientId + "-" + a.UserId + "-" + a.AppName
		delete(pa.names, key)
		delete(pa.ports, a.MappingPort)
	}
	return allocs
}

/**
 * LookupPort finds allocated port for client/app
 * @param {string} clientId - Client identifier
 * @param {string} appName - Application name
 * @returns {int, error} Port number or error if not found
 */
func (pa *PortAllocator) LookupPort(clientId, userId, appName string) (*PortAllocation, error) {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	key := clientId + "-" + userId + "-" + appName
	if alloc, exists := pa.names[key]; exists {
		if alloc.Status == Freed {
			return nil, errors.New("port mapping not found")
		}
		return alloc, nil
	}
	return nil, errors.New("port mapping not found")
}

func (pa *PortAllocator) QueryPorts(clientId, userId, appName string) []PortAllocation {
	ports := []PortAllocation{}

	pa.mu.Lock()
	defer pa.mu.Unlock()

	for _, alloc := range pa.names {
		if clientId != "" && clientId != alloc.ClientId {
			continue
		}
		if userId != "" && userId != alloc.UserId {
			continue
		}
		if appName != "" && appName != alloc.AppName {
			continue
		}
		ports = append(ports, *alloc)
	}
	return ports
}

func (pa *PortAllocator) FreeHangingPorts() []PortAllocation {
	pa.mu.Lock()
	defer pa.mu.Unlock()

	ports := []PortAllocation{}
	now := time.Now()
	for _, alloc := range pa.names {
		if alloc.Status == Allocated && now.Sub(*alloc.AllocTime) > 10*time.Minute {
			alloc.Status = Freed
			ports = append(ports, *alloc)
		}
	}
	for _, a := range ports {
		key := a.ClientId + "-" + a.UserId + "-" + a.AppName
		delete(pa.names, key)
		delete(pa.ports, a.MappingPort)
	}
	return ports
}
