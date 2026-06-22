package configuration

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// 测试用的配置结构体（包装 App）
type testConfig struct {
	App App `json:"app" yaml:"app"`
}

// resetGlobals 重置全局变量到默认值
func resetGlobals() {
	configMu.Lock()
	defer configMu.Unlock()
	defaultConfigPath = "./"
	defaultConfigFileName = "bootstrap"
	defaultConfigFileSuffix = []string{"yml", "yaml", "toml", "json"}
	cmdLinePrefix = "-G"
	cmdLineSplit = "="
}

// ======================== Setter / Getter 测试 ========================

func TestSetDefaultConfigPath(t *testing.T) {
	resetGlobals()
	SetDefaultConfigPath("/tmp/test")
	if got := getDefaultConfigPath(); got != "/tmp/test" {
		t.Errorf("expected /tmp/test, got %s", got)
	}
}

func TestSetDefaultConfigFileName(t *testing.T) {
	resetGlobals()
	SetDefaultConfigFileName("config")
	if got := getDefaultConfigFileName(); got != "config" {
		t.Errorf("expected config, got %s", got)
	}
}

func TestSetDefaultConfigFileSuffix(t *testing.T) {
	resetGlobals()
	suffixes := []string{"json", "toml"}
	SetDefaultConfigFileSuffix(suffixes)
	got := getDefaultConfigFileSuffix()
	if !reflect.DeepEqual(got, suffixes) {
		t.Errorf("expected %v, got %v", suffixes, got)
	}
}

func TestSetCommandLinePrefix(t *testing.T) {
	resetGlobals()
	SetCommandLinePrefix("--")
	if got := getCmdLinePrefix(); got != "--" {
		t.Errorf("expected --, got %s", got)
	}
}

func TestSetCommandLineSplit(t *testing.T) {
	resetGlobals()
	SetCommandLineSplit(":")
	if got := getCmdLineSplit(); got != ":" {
		t.Errorf("expected :, got %s", got)
	}
}

func TestSetDefaultConfigPath_Concurrent(t *testing.T) {
	resetGlobals()
	var wg sync.WaitGroup
	n := 50
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			SetDefaultConfigPath("/tmp/test" + string(rune('0'+idx%10)))
			_ = getDefaultConfigPath()
		}(i)
	}
	wg.Wait()
	// 无 race detector 报错即通过
}

// ======================== applyDefaultTagValues 测试 ========================

func TestApplyDefaultTagValues_String(t *testing.T) {
	type cfg struct {
		Name string `default:"hello"`
		Env  string `default:"dev"`
	}
	c := &cfg{}
	if err := applyDefaultTagValues(c); err != nil {
		t.Fatal(err)
	}
	if c.Name != "hello" {
		t.Errorf("expected Name='hello', got '%s'", c.Name)
	}
	if c.Env != "dev" {
		t.Errorf("expected Env='dev', got '%s'", c.Env)
	}
}

func TestApplyDefaultTagValues_NonZeroPreserved(t *testing.T) {
	type cfg struct {
		Name string `default:"hello"`
	}
	c := &cfg{Name: "existing"}
	if err := applyDefaultTagValues(c); err != nil {
		t.Fatal(err)
	}
	if c.Name != "existing" {
		t.Errorf("expected existing value preserved, got '%s'", c.Name)
	}
}

func TestApplyDefaultTagValues_Bool(t *testing.T) {
	type cfg struct {
		Enable bool `default:"true"`
	}
	c := &cfg{}
	if err := applyDefaultTagValues(c); err != nil {
		t.Fatal(err)
	}
	if !c.Enable {
		t.Error("expected Enable=true")
	}
}

func TestApplyDefaultTagValues_Int(t *testing.T) {
	type cfg struct {
		Port int `default:"8080"`
	}
	c := &cfg{}
	if err := applyDefaultTagValues(c); err != nil {
		t.Fatal(err)
	}
	if c.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", c.Port)
	}
}

func TestApplyDefaultTagValues_Float(t *testing.T) {
	type cfg struct {
		Rate float64 `default:"3.14"`
	}
	c := &cfg{}
	if err := applyDefaultTagValues(c); err != nil {
		t.Fatal(err)
	}
	if c.Rate != 3.14 {
		t.Errorf("expected Rate=3.14, got %f", c.Rate)
	}
}

func TestApplyDefaultTagValues_NonPointerError(t *testing.T) {
	type cfg struct {
		Name string `default:"x"`
	}
	c := cfg{}
	err := applyDefaultTagValues(c)
	if err == nil {
		t.Error("expected error for non-pointer")
	}
}

func TestApplyDefaultTagValues_NilPointerError(t *testing.T) {
	type cfg struct {
		Name string `default:"x"`
	}
	var c *cfg
	err := applyDefaultTagValues(c)
	if err == nil {
		t.Error("expected error for nil pointer")
	}
}

// ======================== mapToProperties 测试 ========================

func TestMapToProperties_Flat(t *testing.T) {
	m := map[string]any{
		"name": "test",
		"port": 8080,
	}
	props := make(map[string]any)
	mapToProperties(m, props, "")
	if props["name"] != "test" {
		t.Errorf("expected name=test, got %v", props["name"])
	}
	if props["port"] != 8080 {
		t.Errorf("expected port=8080, got %v", props["port"])
	}
}

func TestMapToProperties_Nested(t *testing.T) {
	m := map[string]any{
		"app": map[string]any{
			"name": "myapp",
			"env":  "dev",
		},
	}
	props := make(map[string]any)
	mapToProperties(m, props, "")
	if props["app.name"] != "myapp" {
		t.Errorf("expected app.name=myapp, got %v", props["app.name"])
	}
	if props["app.env"] != "dev" {
		t.Errorf("expected app.env=dev, got %v", props["app.env"])
	}
}

func TestMapToProperties_CaseInsensitive(t *testing.T) {
	m := map[string]any{
		"AppName": "myapp",
	}
	props := make(map[string]any)
	mapToProperties(m, props, "")
	if _, ok := props["appname"]; !ok {
		t.Error("expected lowercase key 'appname'")
	}
}

// ======================== propertiesToMap 测试 ========================

func TestPropertiesToMap_RoundTrip(t *testing.T) {
	original := map[string]any{
		"name": "old",
		"app": map[string]any{
			"env": "old",
		},
	}
	// 1. struct → properties
	props := make(map[string]any)
	mapToProperties(original, props, "")

	// 2. 修改 properties
	props["name"] = "new"
	props["app.env"] = "prod"

	// 3. properties → struct (map)
	propertiesToMap(original, props, "")

	if original["name"] != "new" {
		t.Errorf("expected name=new, got %v", original["name"])
	}
	nested := original["app"].(map[string]any)
	if nested["env"] != "prod" {
		t.Errorf("expected app.env=prod, got %v", nested["env"])
	}
}

func TestPropertiesToMap_EmptyValueSkipped(t *testing.T) {
	original := map[string]any{
		"name": "old",
	}
	props := map[string]any{
		"name": "",
	}
	propertiesToMap(original, props, "")
	if original["name"] != "old" {
		t.Errorf("expected old value preserved for empty string, got %v", original["name"])
	}
}

// ======================== isNonZeroValue 测试 ========================

func TestIsNonZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		v        any
		expected bool
	}{
		{"nil", nil, false},
		{"empty string", "", false},
		{"non-empty string", "hello", true},
		{"false bool", false, true}, // bool false 也是有效值
		{"true bool", true, true},
		{"zero int", 0, true}, // 0 也是有效的整数值
		{"positive int", 42, true},
		{"negative int", -1, true},
		{"zero float", 0.0, true},
		{"empty slice", []string{}, false},
		{"non-empty slice", []string{"a"}, true},
		{"empty map", map[string]any{}, false},
		{"non-empty map", map[string]any{"k": "v"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNonZeroValue(tt.v); got != tt.expected {
				t.Errorf("isNonZeroValue(%v) = %v, want %v", tt.v, got, tt.expected)
			}
		})
	}
}

// ======================== fileExists 测试 ========================

func TestFileExists_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	f, err := os.Create(filepath.Join(tmpDir, "test.yml"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	if !fileExists(filepath.Join(tmpDir, "test.yml")) {
		t.Error("expected file to exist")
	}
}

func TestFileExists_NonExistingFile(t *testing.T) {
	if fileExists("/nonexistent/path/file.yml") {
		t.Error("expected file to not exist")
	}
}

func TestFileExists_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	if fileExists(tmpDir) {
		t.Error("expected directory to return false")
	}
}

// ======================== getConfigurationFile 测试 ========================

func TestGetConfigurationFile_Found(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	// 创建 bootstrap.yml
	f, err := os.Create(filepath.Join(tmpDir, "bootstrap.yml"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	SetDefaultConfigPath(tmpDir + "/")
	SetDefaultConfigFileName("bootstrap")
	SetDefaultConfigFileSuffix([]string{"yml"})

	cf := getConfigurationFile()
	if cf == nil {
		t.Fatal("expected configuration file to be found")
	}
	if cf.Suffix != "yml" {
		t.Errorf("expected suffix yml, got %s", cf.Suffix)
	}
	if cf.FileName != "bootstrap" {
		t.Errorf("expected FileName bootstrap, got %s", cf.FileName)
	}
}

func TestGetConfigurationFile_NotFound(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	SetDefaultConfigPath(tmpDir + "/")
	SetDefaultConfigFileName("nonexistent")
	SetDefaultConfigFileSuffix([]string{"yml"})

	cf := getConfigurationFile()
	if cf != nil {
		t.Error("expected nil for non-existent config file")
	}
}

func TestGetConfigurationFile_SecondSuffix(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	// 创建 bootstrap.yaml（第二个后缀）
	f, err := os.Create(filepath.Join(tmpDir, "bootstrap.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	SetDefaultConfigPath(tmpDir + "/")
	SetDefaultConfigFileName("bootstrap")
	SetDefaultConfigFileSuffix([]string{"yml", "yaml"})

	cf := getConfigurationFile()
	if cf == nil {
		t.Fatal("expected configuration file to be found via second suffix")
	}
	if cf.Suffix != "yaml" {
		t.Errorf("expected suffix yaml, got %s", cf.Suffix)
	}
}

// ======================== LoadFile 测试 ========================

func TestLoadFile_ValidFilename(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	// 创建配置文件
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), "name: testapp")

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "bootstrap.yml")
	if err != nil {
		t.Fatal(err)
	}
	// 验证文件名解析正确
	if getDefaultConfigFileName() != "bootstrap" {
		t.Errorf("expected FileName=bootstrap, got %s", getDefaultConfigFileName())
	}
	if suffix := getDefaultConfigFileSuffix(); len(suffix) != 1 || suffix[0] != "yml" {
		t.Errorf("expected Suffix=[yml], got %v", suffix)
	}
}

func TestLoadFile_OnlyPath(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), "app:\n  name: pathonly")

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.App.Name != "pathonly" {
		t.Errorf("expected name=pathonly, got %s", c.App.Name)
	}
}

func TestLoadFile_OnlyFilename(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	// 需要先设置路径
	SetDefaultConfigPath(tmpDir + "/")
	writeFile(t, filepath.Join(tmpDir, "config.json"), `{"app":{"name":"jsonapp"}}`)

	c := &testConfig{}
	err := LoadFile(c, "", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	if c.App.Name != "jsonapp" {
		t.Errorf("expected name=jsonapp, got %s", c.App.Name)
	}
	if getDefaultConfigFileName() != "config" {
		t.Errorf("expected FileName=config, got %s", getDefaultConfigFileName())
	}
}

func TestLoadFile_InvalidFilename_NoDot(t *testing.T) {
	resetGlobals()
	c := &testConfig{}
	err := LoadFile(c, "/tmp", "noextension")
	if err == nil {
		t.Error("expected error for filename without extension")
	}
}

func TestLoadFile_InvalidFilename_EndsWithDot(t *testing.T) {
	resetGlobals()
	c := &testConfig{}
	err := LoadFile(c, "/tmp", "bootstrap.")
	if err == nil {
		t.Error("expected error for filename ending with dot")
	}
}

func TestLoadFile_InvalidFilename_DotOnly(t *testing.T) {
	resetGlobals()
	c := &testConfig{}
	err := LoadFile(c, "/tmp", ".yml")
	if err == nil {
		t.Error("expected error for filename starting with dot")
	}
}

// ======================== Load 集成测试 ========================

func TestLoad_DefaultTags(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	SetDefaultConfigPath(tmpDir + "/")

	c := &testConfig{}
	err := Load(c)
	if err != nil {
		t.Fatal(err)
	}
	// App.Env 有 default:"dev" tag
	if c.App.Env != "dev" {
		t.Errorf("expected Env=dev from default tag, got %s", c.App.Env)
	}
	// App.Name 有 default:"Gin Application" tag
	if c.App.Name != "Gin Application" {
		t.Errorf("expected Name='Gin Application' from default tag, got '%s'", c.App.Name)
	}
}

func TestLoad_FromYAMLFile(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	yamlContent := `
app:
  name: yamltest
  env: production
  web:
    addr: ":9090"
`
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), yamlContent)

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.App.Name != "yamltest" {
		t.Errorf("expected name=yamltest, got %s", c.App.Name)
	}
	if c.App.Env != "production" {
		t.Errorf("expected env=production, got %s", c.App.Env)
	}
	if c.App.Web.Addr != ":9090" {
		t.Errorf("expected addr=:9090, got %s", c.App.Web.Addr)
	}
}

func TestLoad_FromJSONFile(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	jsonContent := `{"app":{"name":"jsontest","env":"staging"}}`
	writeFile(t, filepath.Join(tmpDir, "bootstrap.json"), jsonContent)

	SetDefaultConfigFileSuffix([]string{"json"})
	c := &testConfig{}
	err := LoadFile(c, tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.App.Name != "jsontest" {
		t.Errorf("expected name=jsontest, got %s", c.App.Name)
	}
	if c.App.Env != "staging" {
		t.Errorf("expected env=staging, got %s", c.App.Env)
	}
}

func TestLoad_FromTOMLFile(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	tomlContent := `
[app]
name = "tomltest"
env = "testing"
`
	writeFile(t, filepath.Join(tmpDir, "bootstrap.toml"), tomlContent)

	SetDefaultConfigFileSuffix([]string{"toml"})
	c := &testConfig{}
	err := LoadFile(c, tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	if c.App.Name != "tomltest" {
		t.Errorf("expected name=tomltest, got %s", c.App.Name)
	}
}

func TestLoad_NoConfigFile_NoError(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	SetDefaultConfigPath(tmpDir + "/")

	c := &testConfig{}
	err := Load(c)
	if err != nil {
		t.Fatalf("expected no error when no config file exists, got: %v", err)
	}
	// 应该使用 default tag 值
	if c.App.Name != "Gin Application" {
		t.Errorf("expected default name, got %s", c.App.Name)
	}
}

// ======================== loadEnv 测试 ========================

func TestLoad_EnvVariable(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	SetDefaultConfigPath(tmpDir + "/")

	// 设置环境变量 APP_NAME
	t.Setenv("APP_NAME", "envapp")

	c := &testConfig{}
	err := Load(c)
	if err != nil {
		t.Fatal(err)
	}
	// 环境变量应该优先于 default tag
	if c.App.Name != "envapp" {
		t.Errorf("expected name=envapp from env, got %s", c.App.Name)
	}
}

func TestLoad_EnvVariableLowerPriorityThanFile(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	yamlContent := `
app:
  name: fileapp
`
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), yamlContent)

	t.Setenv("APP_NAME", "envapp")

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "")
	if err != nil {
		t.Fatal(err)
	}
	// 文件优先级高于环境变量
	if c.App.Name != "fileapp" {
		t.Errorf("expected name=fileapp (file > env), got %s", c.App.Name)
	}
}

// ======================== loadFile 环境特定配置测试 ========================

func TestLoad_EnvSpecificFile(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()

	// 默认配置文件，env=dev
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), `
app:
  name: baseapp
  env: dev
`)
	// 环境特定配置文件（优先级更高）
	writeFile(t, filepath.Join(tmpDir, "bootstrap-dev.yml"), `
app:
  name: devapp
`)

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "bootstrap.yml")
	if err != nil {
		t.Fatal(err)
	}
	// 环境特定配置应覆盖默认配置
	if c.App.Name != "devapp" {
		t.Errorf("expected name=devapp from env-specific file, got %s", c.App.Name)
	}
}

func TestLoad_EnvSpecificFileNotExists(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()

	// 只有默认配置文件，env=production
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), `
app:
  name: baseapp
  env: production
`)
	// 不创建 bootstrap-production.yml

	c := &testConfig{}
	err := LoadFile(c, tmpDir, "bootstrap.yml")
	if err != nil {
		t.Fatal(err)
	}
	// 环境特定配置文件不存在时，使用默认配置
	if c.App.Name != "baseapp" {
		t.Errorf("expected name=baseapp, got %s", c.App.Name)
	}
}

// ======================== loadCommandLine 测试 ========================

func TestLoadCommandLine_KeyValue(t *testing.T) {
	resetGlobals()
	saveArgs := os.Args
	defer func() { os.Args = saveArgs }()

	os.Args = []string{"app", "-Gapp.name=cmdapp"}

	props := map[string]any{
		"app.name": "",
	}
	loadCommandLine(props)

	if props["app.name"] != "cmdapp" {
		t.Errorf("expected app.name=cmdapp, got %v", props["app.name"])
	}
}

func TestLoadCommandLine_KeySpaceValue(t *testing.T) {
	resetGlobals()
	saveArgs := os.Args
	defer func() { os.Args = saveArgs }()

	os.Args = []string{"app", "-Gapp.name", "spacedapp"}

	props := map[string]any{
		"app.name": "",
	}
	loadCommandLine(props)

	if props["app.name"] != "spacedapp" {
		t.Errorf("expected app.name=spacedapp, got %v", props["app.name"])
	}
}

func TestLoadCommandLine_UnknownKey(t *testing.T) {
	resetGlobals()
	saveArgs := os.Args
	defer func() { os.Args = saveArgs }()

	os.Args = []string{"app", "-Gunknown.key=value"}

	props := map[string]any{
		"app.name": "",
	}
	loadCommandLine(props) // 不应 panic，仅 log 警告

	if _, ok := props["unknown.key"]; ok {
		t.Error("unknown key should not be added to properties")
	}
}

func TestLoadCommandLine_CustomPrefix(t *testing.T) {
	resetGlobals()
	saveArgs := os.Args
	defer func() { os.Args = saveArgs }()

	SetCommandLinePrefix("--")
	SetCommandLineSplit("=")
	os.Args = []string{"app", "--app.name=custom"}

	props := map[string]any{
		"app.name": "",
	}
	loadCommandLine(props)

	if props["app.name"] != "custom" {
		t.Errorf("expected app.name=custom, got %v", props["app.name"])
	}
}

// ======================== getThisEnv 测试 ========================

func TestGetThisEnv_FromProperties(t *testing.T) {
	props := map[string]any{
		"app.env": "production",
	}
	env := getThisEnv(props, &testConfig{})
	if env != "production" {
		t.Errorf("expected production, got %s", env)
	}
}

func TestGetThisEnv_FromReflection(t *testing.T) {
	// properties 中没有 app.env 的 key，但 struct 中有 App 类型字段
	props := map[string]any{
		// App.Env 的 property key 由 mapToProperties 生成
	}
	// 注意：mapToProperties 会将 struct 展开，App 类型的字段会被嵌套处理
	// 这里直接测试反射路径：set properties key 为 "app.env"
	props["app.env"] = "staging"

	env := getThisEnv(props, &testConfig{})
	if env != "staging" {
		t.Errorf("expected staging, got %s", env)
	}
}

func TestGetThisEnv_Empty(t *testing.T) {
	props := map[string]any{}
	env := getThisEnv(props, &testConfig{})
	if env != "" {
		t.Errorf("expected empty, got %s", env)
	}
}

func TestGetThisEnv_NilConfig(t *testing.T) {
	props := map[string]any{}
	env := getThisEnv(props, nil)
	if env != "" {
		t.Errorf("expected empty for nil config, got %s", env)
	}
}

// ======================== readConfigFile 测试 ========================

func TestReadConfigFile_ValidYAML(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	yamlContent := `
name: readtest
env: dev
`
	writeFile(t, filepath.Join(tmpDir, "test.yml"), yamlContent)

	cf := &ConfigurationFileStruct{
		Path:     tmpDir + "/",
		FileName: "test",
		Suffix:   "yml",
		FullName: filepath.Join(tmpDir, "test.yml"),
	}
	props := map[string]any{
		"name": "",
		"env":  "",
	}
	err := readConfigFile(props, cf)
	if err != nil {
		t.Fatal(err)
	}
	if props["name"] != "readtest" {
		t.Errorf("expected name=readtest, got %v", props["name"])
	}
	if props["env"] != "dev" {
		t.Errorf("expected env=dev, got %v", props["env"])
	}
}

func TestReadConfigFile_FileNotFound(t *testing.T) {
	resetGlobals()
	cf := &ConfigurationFileStruct{
		Path:     "/nonexistent/",
		FileName: "nofile",
		Suffix:   "yml",
		FullName: "/nonexistent/nofile.yml",
	}
	props := map[string]any{}
	err := readConfigFile(props, cf)
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestReadConfigFile_UnknownKeysNotAdded(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	writeFile(t, filepath.Join(tmpDir, "test.yml"), "name: x\nunknown: y")

	cf := &ConfigurationFileStruct{
		Path:     tmpDir + "/",
		FileName: "test",
		Suffix:   "yml",
		FullName: filepath.Join(tmpDir, "test.yml"),
	}
	props := map[string]any{
		"name": "",
	}
	err := readConfigFile(props, cf)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := props["unknown"]; ok {
		t.Error("unknown key should not be added to properties")
	}
}

// ======================== Load 优先级集成测试 ========================

func TestLoad_Priority_CommandLine_Highest(t *testing.T) {
	resetGlobals()
	saveArgs := os.Args
	defer func() { os.Args = saveArgs }()

	tmpDir := t.TempDir()
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), `
app:
  name: fileapp
`)
	t.Setenv("APP_NAME", "envapp")
	os.Args = []string{"app", "-Gapp.name=cmdapp"}

	c := &testConfig{}
	// 已设置文件路径，注入环境变量，命令行参数
	_ = LoadFile(c, tmpDir, "bootstrap.yml")

	// 命令行 > 文件 > 环境变量
	if c.App.Name != "cmdapp" {
		t.Errorf("expected cmdapp (highest priority), got %s", c.App.Name)
	}
}

func TestLoad_Priority_File_Higher_Than_Env(t *testing.T) {
	resetGlobals()
	tmpDir := t.TempDir()
	writeFile(t, filepath.Join(tmpDir, "bootstrap.yml"), `
app:
  name: fileapp
`)
	t.Setenv("APP_NAME", "envapp")

	c := &testConfig{}
	_ = LoadFile(c, tmpDir, "bootstrap.yml")

	if c.App.Name != "fileapp" {
		t.Errorf("expected fileapp (file > env), got %s", c.App.Name)
	}
}

// ======================== 解析函数测试 ========================

func TestParseBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
		wantErr  bool
	}{
		{"true", true, false},
		{"True", true, false},
		{"TRUE", true, false},
		{"1", true, false},
		{"yes", true, false},
		{"on", true, false},
		{"false", false, false},
		{"0", false, false},
		{"no", false, false},
		{"off", false, false},
		{"invalid", false, true},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseBool(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseBool(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.expected {
				t.Errorf("parseBool(%q) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"8080", 8080},
		{"-42", -42},
		{"0", 0},
		{"  100  ", 100},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := parseInt(tt.input)
			if err != nil {
				t.Errorf("parseInt(%q) unexpected error: %v", tt.input, err)
			}
			if got != tt.expected {
				t.Errorf("parseInt(%q) = %d, want %d", tt.input, got, tt.expected)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	got, err := parseFloat("3.14")
	if err != nil {
		t.Fatal(err)
	}
	if got != 3.14 {
		t.Errorf("expected 3.14, got %f", got)
	}
}

func TestParseUint(t *testing.T) {
	got, err := parseUint("8080")
	if err != nil {
		t.Fatal(err)
	}
	if got != 8080 {
		t.Errorf("expected 8080, got %d", got)
	}
}

// ======================== mapToStruct 测试 ========================

func TestMapToStruct_Basic(t *testing.T) {
	type cfg struct {
		Name string
		Port int
	}
	m := map[string]any{
		"Name": "test",
		"Port": 8080,
	}
	c := &cfg{}
	mapToStruct(m, c)
	if c.Name != "test" {
		t.Errorf("expected Name=test, got %s", c.Name)
	}
	if c.Port != 8080 {
		t.Errorf("expected Port=8080, got %d", c.Port)
	}
}

// ======================== App struct default tag 集成测试 ========================

func TestApp_DefaultTags(t *testing.T) {
	app := &App{}
	_ = applyDefaultTagValues(app)
	if app.Name != "Gin Application" {
		t.Errorf("expected Name='Gin Application', got '%s'", app.Name)
	}
	if app.Env != "dev" {
		t.Errorf("expected Env='dev', got '%s'", app.Env)
	}
}

// ======================== 辅助函数 ========================

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(strings.TrimSpace(content)), 0644); err != nil {
		t.Fatal(err)
	}
}
