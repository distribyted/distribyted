package distribyted

import (
	"context"
	"fmt"
	"time"

	peerstore "github.com/libp2p/go-libp2p-peerstore"

	coreiface "github.com/ipfs/interface-go-ipfs-core"
	path "github.com/ipfs/interface-go-ipfs-core/path"
)

type PstoreMngr struct {
	ctx       context.Context
	dht       coreiface.DhtAPI
	swarm     coreiface.SwarmAPI
	store     peerstore.Peerstore
	rootLinks []path.Path
}

func NewPstoreMngr(
	ctx context.Context,
	rootLinks []path.Path,
	dht coreiface.DhtAPI,
	swarm coreiface.SwarmAPI,
	store peerstore.Peerstore,
) *PstoreMngr {
	return &PstoreMngr{ctx, dht, swarm, store, rootLinks}
}

func (p *PstoreMngr) Start() {
	p.loadKnownRootPeers()
	p.startProvidersFinder()
}

func (p *PstoreMngr) loadKnownRootPeers() {
	for _, peer := range p.store.Peers() {
		pp := peer
		go func() {
			fmt.Println("loading known peer", pp)
			addrs := p.store.Addrs(pp)

			if err := p.swarm.Connect(p.ctx, peerstore.PeerInfo{
				ID:    pp,
				Addrs: addrs,
			}); err != nil {
				fmt.Println(err)
			}
		}()
	}
}

func (p *PstoreMngr) startProvidersFinder() {
	for _, link := range p.rootLinks {
		fmt.Println("loading new root link", link)
		ll := link
		go func() {
			pInfo, err := p.dht.FindProviders(p.ctx, ll)
			if err != nil {
				fmt.Println(err)
				return
			}
			for i := range pInfo {
				fmt.Println("Adding new address from peer", i.ID, "LINK", ll)
				p.store.SetAddrs(i.ID, i.Addrs, 24*30*time.Hour)
			}
		}()
	}
}
