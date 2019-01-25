package notaryapi

import (
	"encoding/json"
	"testing"

	"fmt"

	"sync"

	"github.com/SmartMeshFoundation/distributed-notary/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewKeyGenerationPhase1MessageRequest(t *testing.T) {
	str := "{\"name\":\"PKN-Phase1PubKeyProof\",\"session_id\":\"0x7a6b6242aa8a4bdc4b2814fe50826d41eb1f0318f90c151812c7f2e150f24d75\",\"sender_notary_id\":0,\"signer\":\"0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa\",\"signature\":\"gq0IxEiDcUX77fPdEKdx+5XZO3d9jxx0NLYlEaDh54tBbaGWsN53tr6m64rBAzJ9X7ML9rVlFPZxXhCDcYp2Uxw=\",\"request_id\":\"259c\",\"msg\":{\"Proof\":{\"PK\":{\"X\":10174871828079934626442637541222929157029255613022845095011872425499772812956,\"Y\":5977319075155931208570142654655228345323205805632938290013058988051096876177},\"PkTRandCommitment\":{\"X\":70357495813316055203731525500804010256611375431203492434767353295883130132042,\"Y\":102638225753872609010964689353648519369108941952694072041904845868405090667849},\"ChallengeResponse\":{\"D\":43283280515197778264155963032047148587556394365100441106435252494626102286698}}}}"
	req := &KeyGenerationPhase1MessageRequest{}
	err := json.Unmarshal([]byte(str), &req)
	assert.Empty(t, err)
	fmt.Println(utils.ToJSONStringFormat(req))
}

func TestSyncMap(t *testing.T) {
	m := new(sync.Map)
	q, ok := m.Load(1)
	fmt.Println(q, ok)
	m.Store(1, 11)
	m.Store(2, 22)
	q, ok = m.Load(1)
	fmt.Println(q, ok)
	q, ok = m.Load(2)
	fmt.Println(q, ok)
}
