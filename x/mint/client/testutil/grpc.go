package testutil

import (
	"fmt"

	"github.com/line/lbm-sdk/testutil"
	sdk "github.com/line/lbm-sdk/types"
	grpctypes "github.com/line/lbm-sdk/types/grpc"

	"github.com/gogo/protobuf/proto"

	minttypes "github.com/line/lbm-sdk/x/mint/types"
)

func (s *IntegrationTestSuite) TestQueryGRPC() {
	val := s.network.Validators[0]
	baseURL := val.APIAddress
	testCases := []struct {
		name     string
		url      string
		headers  map[string]string
		respType proto.Message
		expected proto.Message
	}{
		{
			"gRPC request params",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/params", baseURL),
			map[string]string{},
			&minttypes.QueryParamsResponse{},
			&minttypes.QueryParamsResponse{
				Params: minttypes.NewParams("stake", sdk.NewDecWithPrec(13, 2), sdk.NewDecWithPrec(100, 2),
					sdk.NewDec(1), sdk.NewDecWithPrec(67, 2), 60*60*8766/5),
			},
		},
		{
			"gRPC request inflation",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/inflation", baseURL),
			map[string]string{},
			&minttypes.QueryInflationResponse{},
			&minttypes.QueryInflationResponse{
				Inflation: sdk.NewDec(1),
			},
		},
		{
			"gRPC request annual provisions",
			fmt.Sprintf("%s/cosmos/mint/v1beta1/annual_provisions", baseURL),
			map[string]string{
				grpctypes.GRPCBlockHeightHeader: "1",
			},
			&minttypes.QueryAnnualProvisionsResponse{},
			&minttypes.QueryAnnualProvisionsResponse{
				AnnualProvisions: sdk.NewDec(500000000),
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		resp, err := testutil.GetRequestWithHeaders(tc.url, tc.headers)
		s.Run(tc.name, func() {
			s.Require().NoError(err)
			s.Require().NoError(val.ClientCtx.Codec.UnmarshalJSON(resp, tc.respType))
			s.Require().Equal(tc.expected.String(), tc.respType.String())
		})
	}
}
