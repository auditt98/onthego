package generators

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/auditt98/onthego/utils"
)

func InitGenerator() error {
	config, _ := utils.LoadConf()
	for version, versionConfig := range config.Handler.Versions {
		GenerateAPIVersionFolder(config.Handler.Path, version)
		GenerateAPIVersionRouter()
		for groupPath, routeGroup := range versionConfig.Routes {
			path := utils.JoinPaths(config.Handler.Path, version)
			absPath := filepath.Join(path, groupPath+".go")

			// Check if the file already exists
			_, err := os.Stat(absPath)
			if err == nil {
				fmt.Printf("File '%s' already exists\n", absPath)
				continue
			}
			GenerateHandlerFile(path, groupPath, version, "./generators/templates/handler.tmpl")
			// utils.FmtFile(utils.JoinPaths(path, groupPath+".go"))
			GenerateHandlerRouter()
			handlerPath := utils.JoinPaths(path, groupPath+".go")
			GenerateMultipleAction(routeGroup.Actions, groupPath, version, handlerPath, "./generators/templates/action.tmpl")
		}
	}
	return nil
}
