package database

import (
	"fmt"
)

type Script struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	Platform string `json:"platform"`
	Type     string `json:"type"`
	Agent    string `json:"agent"`
	User     string `json:"user"`
}

var ScriptBucketName = []byte("scripts")

func NewScript(scriptInfo map[string]string, active bool) (*Script, error) {
	fields := []string{"name", "code", "platform", "cronschedule"}

	for _, field := range fields {
		value, ok := scriptInfo[field]
		if !ok {
			return nil, fmt.Errorf("field %s not found", value)
		}
	}

	/*
		CronSchedule := scriptInfo["cronschedule"]

		i, err := strconv.ParseInt(CronSchedule, 10, 64)
		if err != nil {
			i = 0
		}
	*/

	return &Script{
		Name:     scriptInfo["name"],
		Code:     scriptInfo["code"],
		Platform: scriptInfo["platform"],
		//CronSchedule: i,
	}, nil
}

func CreateScript(script *Script) error {
	scriptbyte, err := Encode(script)
	if err != nil {
		return err
	}

	err = DB.Create([]byte(script.Name), scriptbyte, ScriptBucketName)
	if err != nil {
		return err
	}
	return nil
}

func UpdateScript(script *Script) error {
	scriptbyte, err := Encode(script)
	if err != nil {
		return err
	}

	err = DB.Update([]byte(script.Name), scriptbyte, ScriptBucketName)
	if err != nil {
		return err
	}
	return nil
}

func GetAllScripts() ([]*Script, error) {
	scripts := []*Script{}

	scriptbyte, err := DB.ReadAll(ScriptBucketName)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	for _, val := range scriptbyte {
		script := &Script{}
		err = Decode(val, script)
		scripts = append(scripts, script)

		if err != nil {
			return nil, err
		}

	}
	return scripts, nil
}

func GetScript(name string) (*Script, error) {
	script := &Script{}
	scriptbyte, err := DB.Read([]byte(name), ScriptBucketName)

	if err != nil {
		return nil, err
	}

	err = Decode(scriptbyte, script)
	if err != nil {
		return nil, err
	}

	return script, nil
}

func DeleteScript(name string) error {

	return DB.Delete([]byte(name), ScriptBucketName)
}
