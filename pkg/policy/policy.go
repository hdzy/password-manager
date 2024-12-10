package policy

import (
	"bytes"
	"encoding/gob"
	"errors"
	"os"
	"path/filepath"
	"reflect"
)

const (
	filePath = "files/policies.dat"
)

var (
	ErrEmpty     = errors.New("no policies in struct")
	ErrNotFound  = errors.New("policy not found")
	ErrFieldName = errors.New("no such field in struct")
	ErrValueType = errors.New("not excepted type of value for this field")
)

type Policy struct {
	Name           string `json:"name" print:"Название:"`
	Symbols        []rune `json:"symbols" print:"Доступные символы:"`
	MinimumNumbers int    `json:"minimum_numbers" print:"Минимальное количество цифр:"`
	MinTopRegister int    `json:"min_top_register" print:"Минимальное количество верхнего регистра:"`
	MinBotRegister int    `json:"min_bot_register" print:"Минимальное количество нижнего регистра:"`
	MinSpec        int    `json:"min_spec" print:"Минимальное количество спецсимволов:"`
	MinNumbers     int    `json:"min_numbers" print:"Минимальное количество цифр:"`
	CharProc       int    `json:"char_proc" print:"Процент символов:"`
	NumbProc       int    `json:"numb_proc" print:"Процент цифр:"`
	SpecProc       int    `json:"spec_proc" print:"Процент спецсимволов:"`
}

// Fvm - FieldValueMap for change fields in policy struct
type Fvm map[string]interface{}

func NewPolicy(name string, symbols []rune, minimumNumbers, minTopRegister, minBotRegister, minSpec, minNumbers, charProc, numbProc, specProc int) *Policy {
	return &Policy{
		Name:           name,
		Symbols:        symbols,
		MinimumNumbers: minimumNumbers,
		MinTopRegister: minTopRegister,
		MinBotRegister: minBotRegister,
		MinSpec:        minSpec,
		MinNumbers:     minNumbers,
		CharProc:       charProc,
		NumbProc:       numbProc,
		SpecProc:       specProc,
	}
}

// update values {v} of fields {f} where fvm = map[f]v
func (p *Policy) update(fvm Fvm) error {
	v := reflect.ValueOf(p)

	for key, value := range fvm {
		field := v.Elem().FieldByName(key)

		if !field.IsValid() {
			return ErrFieldName
		}

		val := reflect.ValueOf(value)

		if field.Type() != val.Type() {
			return ErrValueType
		}

		field.Set(val) // Изменяем поле
	}

	return nil
}

type Policies []*Policy

// openFile Open file policies.json with flags
func (p *Policies) openFile(flags int) (*os.File, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(projectDir, filePath)

	file, err := os.OpenFile(fullPath, flags, 0644)
	if os.IsNotExist(err) {
		file, err = os.Create(fullPath)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	return file, nil
}

// Load JSON data from file
func (p *Policies) Load() error {
	file, err := p.openFile(os.O_RDONLY)
	if err != nil {
		return err
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(file)
	if err != nil {
		return err
	}

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(p)
	if err != nil {
		return err
	}

	return nil
}

// Save JSON data to file
func (p *Policies) Save() error {
	file, err := p.openFile(os.O_WRONLY | os.O_TRUNC)
	if err != nil {
		return err
	}
	defer file.Close()

	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err = encoder.Encode(p)
	if err != nil {
		return err
	}

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// UpdateByName Find first policy with Policy.Name and change fields got from Fvm
func (p *Policies) UpdateByName(name string, fvm Fvm) error {
	if len(*p) == 0 {
		return ErrEmpty
	}

	for _, policy := range *p {
		if policy.Name == name {
			err := policy.update(fvm)
			if err != nil {
				return err
			}
			err = p.Save()
			return err
		}
	}

	return ErrNotFound
}

// RemoveByName Find first policy with name and removes it
func (p *Policies) RemoveByName(name string) error {
	if len(*p) == 0 {
		return ErrEmpty
	}

	for i, policy := range *p {
		if policy.Name == name {
			*p = append((*p)[:i], (*p)[i+1:]...)

			return nil
		}
	}

	return ErrNotFound
}

// New push new Policy to Policies
func (p *Policies) New(newEl *Policy) {
	*p = append(*p, newEl)
}
