package utils

import (
	"cmp"
	"encoding/binary"
	"net"
	"sort"

	"github.com/chzyer/logex"
)

func IP2Int(ip net.IP) uint32 {
	return binary.BigEndian.Uint32([]byte(ip))
}

func Int2IP(ip uint32) net.IP {
	newIP := make([]byte, 4)
	binary.BigEndian.PutUint32(newIP, ip)
	return newIP
}

var MASK_LOOKUP = func() map[int]net.IPMask {
	out := make(map[int]net.IPMask)
	for i := 0; i < 32; i++ {
		val := (1 << i) - 1
		out[val] = net.CIDRMask(32-i, 32)
	}
	return out
}()

func Range2CIDR(start, end uint32) *net.IPNet {
	ip := Int2IP(start)
	delta := (end - start)
	mask, ok := MASK_LOOKUP[int(delta)]
	if !ok {
		panic("invalid delta")
	}

	return &net.IPNet{
		IP:   ip,
		Mask: mask,
	}
}

func CIDR2Range(cidr string) (uint32, uint32, error) {
	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return 0, 0, logex.Trace(err, cidr)
	}
	ones, bits := ipnet.Mask.Size()
	delta := uint32(1<<(bits-ones) - 1)
	start := IP2Int(ipnet.IP)
	end := start + delta

	return start, end, nil
}

func MapKeys[T any](m map[string]T) []string {
	n := make([]string, 0, len(m))
	for k := range m {
		n = append(n, k)
	}
	sort.Slice(n, func(i, j int) bool {
		return n[i] < n[j]
	})
	return n
}

func CollectSliceValues[T any](m []T, c func(T) string) []string {
	n := make([]string, 0, len(m))
	for _, v := range m {
		n = append(n, c(v))
	}
	sort.Slice(n, func(i, j int) bool {
		return n[i] < n[j]
	})
	return n
}

func CollectMapValues[K comparable, T any](m map[K]T, c func(T) string) []string {
	n := make([]string, 0, len(m))
	for _, v := range m {
		n = append(n, c(v))
	}
	sort.Slice(n, func(i, j int) bool {
		return n[i] < n[j]
	})
	return n
}

func SortUniqueSlice[T cmp.Ordered](lib []T) []T {
	newLib := make([]T, 0, len(lib))
	m := make(map[T]int, len(lib))
	for i := 0; i < len(lib); i++ {
		if _, ok := m[lib[i]]; !ok {
			newLib = append(newLib, lib[i])
			m[lib[i]] = i
		}
	}
	sort.Slice(newLib, func(i, j int) bool {
		return newLib[i] < newLib[j]
	})
	return newLib
}

func SliceToMapIndex(lib []string) map[string]int {
	m := make(map[string]int, len(lib))
	for i := 0; i < len(lib); i++ {
		m[lib[i]] = i
	}
	return m
}

type IpRangeRecord interface {
	GetStartIP() uint32
	GetEndIP() uint32
}

func FindIpInfo[T IpRangeRecord](lib []T, ip net.IP) (T, bool) {
	ipInt := IP2Int(ip)
	idx, found := sort.Find(len(lib), func(i int) int {
		ipRange := lib[i]
		start := ipRange.GetStartIP()
		end := ipRange.GetEndIP()
		if start <= ipInt && ipInt <= end {
			return 0
		} else if ipInt < start {
			return -1
		} else {
			return 1
		}
	})
	if !found {
		var n T
		return n, false
	}
	return lib[idx], true
}
