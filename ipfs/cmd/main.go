package main

import (
	"context"
	"log"
	"time"

	"github.com/ajnavarro/distribyted"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/ipfs/go-cid"
	coreiface "github.com/ipfs/interface-go-ipfs-core"
	"github.com/ipfs/interface-go-ipfs-core/path"
	pstoreds "github.com/libp2p/go-libp2p-peerstore/pstoreds"
)

const ipnsNameKey = "ipns_test"

func stats(ctx context.Context, api coreiface.CoreAPI) {
	for {
		// do some job
		kaddrs, err := api.Swarm().KnownAddrs(ctx)
		if err != nil {
			log.Println(err)
		}

		log.Println("known address", len(kaddrs))

		peers, err := api.Swarm().Peers(ctx)
		if err != nil {
			log.Println(err)
		}

		log.Println("peers", len(peers))

		time.Sleep(5000 * time.Millisecond)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api, repo, err := distribyted.Spawn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	su, err := repo.GetStorageUsage()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("STORAGE USAGE", su)
	go stats(ctx, api)

	log.Println("Init peer store")
	pstore, err := pstoreds.NewPeerstore(ctx, repo.Datastore(), pstoreds.Options{
		CacheSize:           1024,
		GCPurgeInterval:     1 * time.Minute,
		GCLookaheadInterval: 2 * time.Minute,
		GCInitialDelay:      10 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	psm := distribyted.NewPstoreMngr(
		ctx,
		[]path.Path{path.New("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW")}, //path.New("QmPfyHeZgnf6nEoMKogLGxMeDi4YkwXDi8NbcxFdnqu8kX"),

		api.Dht(),
		api.Swarm(),
		pstore,
	)

	psm.Start()

	id, err := cid.Parse("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW")
	if err != nil {
		log.Fatal(err)
	}

	opts := &fs.Options{}
	opts.Debug = false

	node, err := api.Dag().Get(ctx, id)
	if err != nil {
		log.Fatalf("get node fail: %v\n", err)
	}

	server, err := fs.Mount("/tmp/distribyted", distribyted.NewIPFSRoot(api, node), opts)
	if err != nil {
		log.Fatalf("Mount fail: %v\n", err)
	}

	defer server.Unmount()
	server.Wait()
}

func main3() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api, repo, err := distribyted.Spawn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go stats(ctx, api)

	log.Println("Init peer store")
	pstore, err := pstoreds.NewPeerstore(ctx, repo.Datastore(), pstoreds.Options{
		CacheSize:           1024,
		GCPurgeInterval:     1 * time.Minute,
		GCLookaheadInterval: 2 * time.Minute,
		GCInitialDelay:      10 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	psm := distribyted.NewPstoreMngr(
		ctx,
		[]path.Path{path.New("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW")}, //path.New("QmPfyHeZgnf6nEoMKogLGxMeDi4YkwXDi8NbcxFdnqu8kX"),

		api.Dht(),
		api.Swarm(),
		pstore,
	)

	psm.Start()

	log.Println("Getting object")

	nn, err := api.Object().Get(ctx, path.New("QmdZ5du3xjbTK3HQpYjxHzGMFWJwbMnCWeG3tV7LH8D7Bk"))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("NODE", nn.String())
	for _, li := range nn.Links() {
		log.Println("LINK", li)
		log.Println("LINK NAME", li.Name)
	}

	cID, err := cid.Parse("QmdZ5du3xjbTK3HQpYjxHzGMFWJwbMnCWeG3tV7LH8D7Bk")
	if err != nil {
		log.Fatal(err)
	}

	dagNode, err := api.Dag().Get(ctx, cID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DAG NODE", dagNode.String())
	for _, li := range dagNode.Links() {
		log.Println("DAG LINK", li)
		log.Println("DAG LINK NAME", li.Name)
	}

	stats, err := dagNode.Stat()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("DAG STATS", stats)

	// opts := &fs.Options{}
	// opts.Debug = true
	// server, err := fs.Mount("/tmp/distribyted", distribyted.NewIPFSRoot(api), opts)
	// if err != nil {
	// 	log.Fatalf("Mount fail: %v\n", err)
	// }

	// defer server.Unmount()
	// server.Wait()
}

func main2() {
	log.Println("STARTING TEST")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	api, _, err := distribyted.Spawn(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go stats(ctx, api)
	// // add directory
	// resolved, err := api.Unixfs().Add(ctx,
	// 	files.NewBytesFile([]byte("HELLO ANOTHER ANOTHER WORLD")), options.Unixfs.Pin(true))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("added directory with multihash: ", resolved)

	// key, err := api.Key().Generate(ctx, ipnsNameKey)
	// if err != nil {
	// 	log.Println(err)
	// 	keys, err := api.Key().List(ctx)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	for _, k := range keys {
	// 		log.Println("KEY", k.Name(), k.Path(), k.ID())
	// 		if k.Name() == ipnsNameKey {
	// 			key = k
	// 			break
	// 		}
	// 	}
	// }

	// ipnsPublished, err := api.Name().Publish(ctx, resolved,
	// 	options.Name.AllowOffline(true),
	// 	options.Name.Key(key.Name()),
	// 	options.Name.ValidTime(1000*time.Hour),
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("added a name to the directory: ", ipnsPublished.Name(), ipnsPublished.Value().String())

	// ipnsResolved, err := api.Name().
	// 	Resolve(ctx, "QmPsBX4fu5zgkPVFpQRbEVt2cCs4V9b1gBFjdnhesZMNhY")
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("resolved OLD name!: ", ipnsResolved)

	// ipnsResolved, err = api.Name().Resolve(ctx, ipnsPublished.Name())
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("resolved name!: ", ipnsResolved)

	// str := "{\"Addrs\":[\"/ip4/167.71.52.253/tcp/4001\",\"/ip4/165.227.144.202/tcp/4001\",\"/ip4/167.71.52.253/tcp/4002/ws\",\"/ip6/::1/tcp/4002/ws\",\"/ip4/10.19.0.13/tcp/4002/ws\",\"/ip4/127.0.0.1/tcp/4002/ws\",\"/ip4/165.227.144.202/tcp/4002/ws\",\"/ip6/::1/tcp/4001\",\"/ip4/10.19.0.13/tcp/4001\",\"/ip4/127.0.0.1/tcp/4001\",\"/ip4/10.19.0.9/tcp/4001\",\"/ip6/2a03:b0c0:3:d0::3e3:7001/tcp/4002/ws\",\"/ip6/2a03:b0c0:3:d0::3e3:7001/tcp/4001\",\"/ip4/10.19.0.9/tcp/4002/ws\"],\"ID\":\"QmSRg4CQN4aSTKDahNSjwE2BnMRjZkthS5mdmcnau85FM5\"}"
	// var pp pstore.PeerInfo

	// if err := pp.UnmarshalJSON([]byte(str)); err != nil {
	// 	log.Fatal(err)
	// }
	// if err := api.Swarm().Connect(ctx, pp); err != nil {
	// 	log.Println(err)
	// }

	// peersChan, err := api.Dht().FindProviders(ctx, path.New("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW"))
	// if err != nil {
	// 	log.Println(err)
	// }

	// log.Println("listing peers...")
	// for peer := range peersChan {
	// 	data, err := peer.MarshalJSON()
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// 	log.Println("DATA", string(data))
	// 	log.Println("peer", peer)
	// }

	log.Println("getting object...")

	//nn, err := api.Object().Get(ctx, path.New("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW/snes/2020 Super Baseball (USA).7z"))
	//nn, err := api.Object().Get(ctx, path.New("QmayX17sMGH1dH4sAa16ZgBYSrMUMSjx7qVKU73aqWFf2M"))
	//nn, err := api.Object().Get(ctx, path.New("QmPfyHeZgnf6nEoMKogLGxMeDi4YkwXDi8NbcxFdnqu8kX"))
	//nn, err := api.Object().Get(ctx, path.New("QmP7ESfBXavxftW3GCpbtJcFUdPi65uBe2YTMS53dWprZW"))
	nn, err := api.Object().Get(ctx, path.New("QmdZ5du3xjbTK3HQpYjxHzGMFWJwbMnCWeG3tV7LH8D7Bk"))

	if err != nil {
		log.Fatal(err)
	}

	log.Println("NODE", nn.String())
	for _, li := range nn.Links() {
		log.Println("LINK", li)
	}
}
