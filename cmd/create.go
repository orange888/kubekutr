package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
	models "zerodha.tech/janus/models"
	"zerodha.tech/janus/utils"
)

// CreateResource creates ad hoc resources and outputs manifests.
func (hub *Hub) CreateResource(config models.Config) cli.Command {
	return cli.Command{
		Name:    "create",
		Aliases: []string{"c"},
		Usage:   "Create ad-hoc resource(s) and merge with existing GitOps directory",
		Action:  hub.initApp(hub.create),
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "resource, r",
				Usage: "Create manifests for only a particular resource `TYPE`",
			},
			cli.StringFlag{
				Name:  "name, n",
				Usage: "`NAME` of the resource to lookup in config. Use in combination with --resource",
			},
			cli.StringFlag{
				Name:  "path, p",
				Usage: "`PATH` to an existing GitOps directory generated by scaffold",
			},
			cli.StringFlag{
				Name:  "config",
				Usage: "Path to `FILE` config",
			},
		},
	}
}

func (hub *Hub) create(cliCtx *cli.Context) error {
	var (
		projectDir = utils.GetRootDir(cliCtx.String("path"))
	)
	err := utils.LookupGitopsDirectory(subPaths, projectDir)
	if err != nil {
		hub.Logger.Errorf("Output directory %s doesn't match the expected GitOps directory structure generated from `scaffold`", cliCtx.String("path"))
		return err
	}
	// If a particular resource type is selected, create manifest only for that
	switch cliCtx.String("resource") {
	case "deployment":
		for _, dep := range hub.Config.Deployments {
			if cliCtx.String("name") != "" {
				if cliCtx.String("name") != dep.Name {
					continue
				}
			}
			err := loadDeployment(dep, filepath.Join(projectDir, "base", "deployments"))
			if err != nil {
				return err
			}
			hub.Logger.Debugf("Created manifest for deployment: %s", dep.Name)
		}
	case "service":
		// Create services
		for _, svc := range hub.Config.Services {
			if cliCtx.String("name") != "" {
				if cliCtx.String("name") != svc.Name {
					continue
				}
			}
			err := loadService(svc, filepath.Join(projectDir, "base", "services"))
			if err != nil {
				return err
			}
			hub.Logger.Debugf("Created manifest for service: %s", svc.Name)
		}
	case "ingress":
		// Create ingress
		for _, ing := range hub.Config.Ingresses {
			if cliCtx.String("name") != "" {
				if cliCtx.String("name") != ing.Name {
					continue
				}
			}
			err := loadIngress(ing, filepath.Join(projectDir, "base", "ingresses"))
			if err != nil {
				return err
			}
			hub.Logger.Debugf("Created manifest for ingress: %s", ing.Name)
		}
	case "":
		fmt.Println("ALL. to be implemented")
	default:
		return fmt.Errorf("Invalid resource type %s selected", cliCtx.String("resource"))
	}
	return nil
}
