package win32

import (
	"golang.org/x/sys/windows/registry"
)

const (
	RegistryDefaultName = ""
)

func GetClassRootValue(path string, name string) (string, error) {
	return GetRegistryValue(registry.CLASSES_ROOT, path, name)
}

func SetClassRootValue(path string, name string, value string) error {
	return CreateRegistryValue(registry.CLASSES_ROOT, path, name, value)
}

func GetRegistryValue(regKey registry.Key, path string, name string) (string, error) {
	k, err := registry.OpenKey(regKey, path, registry.READ)
	if err != nil {
		return "", err
	}
	defer k.Close()

	s, _, err := k.GetStringValue(name)
	if err != nil {
		return "", err
	}
	return s, nil
}

func CreateRegistryValue(regKey registry.Key, path string, name string, value string) error {
	k, _, err := registry.CreateKey(regKey, path, registry.WRITE)
	if err != nil {
		return err
	}
	defer k.Close()

	return k.SetStringValue(name, value)
}
