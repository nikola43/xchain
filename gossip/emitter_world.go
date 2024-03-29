package gossip

import (
	"sync/atomic"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/galaxy-foundation/icicb-base/hash"
	"github.com/galaxy-foundation/icicb-base/inter/idx"

	"github.com/goicicb/evmcore"
	"github.com/goicicb/inter"
	"github.com/goicicb/utils/wgmutex"
	"github.com/goicicb/valkeystore"
	"github.com/goicicb/vecmt"
)

// emitterWorld implements emitter.World interface
type emitterWorld struct {
	s *Service

	*Store
	*wgmutex.WgMutex
	*evmcore.TxPool
	valkeystore.SignerI
	types.Signer
}

func (ew *emitterWorld) Check(emitted *inter.EventPayload, parents inter.Events) error {
	// sanity check
	return ew.s.checkers.Validate(emitted, parents.Interfaces())
}

func (ew *emitterWorld) Process(emitted *inter.EventPayload) error {
	return ew.s.processEvent(emitted)
}

func (ew *emitterWorld) Broadcast(emitted *inter.EventPayload) {
	// PM listens and will broadcast it
	ew.s.feed.newEmittedEvent.Send(emitted)
}

func (ew *emitterWorld) Build(e *inter.MutableEventPayload, onIndexed func()) error {
	return ew.s.buildEvent(e, onIndexed)
}

func (ew *emitterWorld) DagIndex() *vecmt.Index {
	return ew.s.dagIndexer
}

func (ew *emitterWorld) IsBusy() bool {
	return atomic.LoadUint32(&ew.s.eventBusyFlag) != 0 || atomic.LoadUint32(&ew.s.blockBusyFlag) != 0
}

func (ew *emitterWorld) IsSynced() bool {
	return atomic.LoadUint32(&ew.s.pm.synced) != 0
}

func (ew *emitterWorld) PeersNum() int {
	return ew.s.pm.peers.Len()
}

func (ew *emitterWorld) GetHeads(epoch idx.Epoch) hash.Events {
	return ew.Store.GetHeadsSlice(epoch)
}

func (ew *emitterWorld) GetLastEvent(epoch idx.Epoch, from idx.ValidatorID) *hash.Event {
	return ew.Store.GetLastEvent(epoch, from)
}

