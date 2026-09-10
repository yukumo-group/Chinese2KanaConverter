package polyphonic

import (
	"encoding/json"
	"maps"
	"os"
	"sync"

	"github.com/yukumo-group/Chinese2KanaConverter/internal/cpyconverter"
)

// Manager manages polyphonic phrases
type Manager struct {
	sync.RWMutex
	Heteronym      map[string]string `json:"heteronym"`
	targetFilePath string
}

// NewManager creates new manager
func NewManager() *Manager {
	return &Manager{
		Heteronym:      make(map[string]string),
		targetFilePath: "polyphonic.json",
	}
}

// NewManagerFromFile creates new manager through reading .json file
func NewManagerFromFile(
	targetFilePath string,
) (*Manager, error) {
	var newManager *Manager
	data, err := os.ReadFile(targetFilePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &newManager)
	return newManager, err
}

// SetTargetFile sets target file for storing polyphonics
func (manager *Manager) SetTargetFile(
	targetFilePath string,
) {
	manager.Lock()
	defer manager.Unlock()
	manager.targetFilePath = targetFilePath
}

// Save saves the file
func (manager *Manager) Save() error {
	manager.Lock()
	defer manager.Unlock()
	data, err := json.Marshal(
		manager,
	)
	if err != nil {
		return err
	}
	err = os.WriteFile(
		manager.targetFilePath,
		data,
		0644,
	)
	return err
}

// load loads the heteronyms.
// **Unlocked!**
func (manager *Manager) load() {
	cpyconverter.DumpHeteronymMap(
		maps.Clone(manager.Heteronym),
	)
}

// Initialize initializes the manager
func (manager *Manager) Initialize() {
	manager.Lock()
	defer manager.Unlock()
	manager.load()
}

// AddPolyphonic adds new polyphonic.
// e.g. "都会区" and "dū huì qū"
func (manager *Manager) AddPolyphonic(
	chinese string,
	pinyin string,
) {
	manager.Lock()
	defer manager.Unlock()
	manager.Heteronym[chinese] = pinyin
	manager.load()
}

// GetData gets polyphonic map
func (manager *Manager) GetData() map[string]string {
	manager.RLock()
	defer manager.RUnlock()
	result := maps.Clone(manager.Heteronym)
	return result
}
