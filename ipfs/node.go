package distribyted

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	config "github.com/ipfs/go-ipfs-config"
	core "github.com/ipfs/go-ipfs/core"
	coreapi "github.com/ipfs/go-ipfs/core/coreapi"
	"github.com/ipfs/go-ipfs/core/node/libp2p"
	"github.com/ipfs/go-ipfs/plugin/loader"
	"github.com/ipfs/go-ipfs/repo"
	coreiface "github.com/ipfs/interface-go-ipfs-core"

	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
)

func setupPlugins(path string) error {
	// Load plugins. This will skip the repo if not available.
	plugins, err := loader.NewPluginLoader(filepath.Join(path, "plugins"))
	if err != nil {
		return fmt.Errorf("error loading plugins: %s", err)
	}

	if err := plugins.Initialize(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	if err := plugins.Inject(); err != nil {
		return fmt.Errorf("error initializing plugins: %s", err)
	}

	return nil
}

func open(ctx context.Context, repoPath string) (coreiface.CoreAPI, repo.Repo, error) {
	// Open the repo
	r, err := fsrepo.Open(repoPath)
	if err != nil {
		return nil, nil, err
	}

	// Construct the node
	node, err := core.NewNode(ctx, &core.BuildCfg{
		Online:    true,
		Permanent: true,
		Routing:   libp2p.DHTOption,
		Repo:      r,
	})
	if err != nil {
		return nil, nil, err
	}

	fmt.Println("IDENTITY", node.Identity)

	api, err := coreapi.NewCoreAPI(node)
	if err != nil {
		return nil, nil, err
	}

	return api, r, nil
}

func tmpNode(ctx context.Context) (coreiface.CoreAPI, repo.Repo, error) {
	// dir, err := ioutil.TempDir("", "ipfs-shell")
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get temp dir: %s", err)
	// }
	dir := "ditribyted"
	cfg, err := config.Init(ioutil.Discard, 2048)
	if err != nil {
		return nil, nil, err
	}

	// configure the temporary node
	cfg.Swarm.ConnMgr.LowWater = 10
	cfg.Swarm.ConnMgr.HighWater = 30
	cfg.Swarm.ConnMgr.GracePeriod = "5s"
	cfg.Swarm.ConnMgr.Type = "basic"

	//cfg.Routing.Type = "dhtclient"
	cfg.Routing.Type = "dht"
	cfg.Experimental.QUIC = true
	cfg.Datastore.Spec = map[string]interface{}{
		"type": "badgerds",
		"path": "badger",
	}

	err = fsrepo.Init(dir, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to init ephemeral node: %s", err)
	}
	return open(ctx, dir)
}

func Spawn(ctx context.Context) (coreiface.CoreAPI, repo.Repo, error) {
	defaultPath, err := config.PathRoot()
	if err != nil {
		// shouldn't be possible
		return nil, nil, err
	}

	if err := setupPlugins(defaultPath); err != nil {
		return nil, nil, err
	}

	return tmpNode(ctx)
}
