package host

import (
	sdk "github.com/line/lbm-sdk/types"
	sdkerrors "github.com/line/lbm-sdk/types/errors"
	capabilitytypes "github.com/line/lbm-sdk/x/capability/types"

	"github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/host/keeper"
	"github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/host/types"
	icatypes "github.com/line/lbm-sdk/x/ibc/applications/27-interchain-accounts/types"
	channeltypes "github.com/line/lbm-sdk/x/ibc/core/04-channel/types"
	ibcexported "github.com/line/lbm-sdk/x/ibc/core/exported"
)

// IBCModule implements the ICS26 interface for interchain accounts host chains
type IBCModule struct {
	keeper keeper.Keeper
}

// NewIBCModule creates a new IBCModule given the associated keeper
func NewIBCModule(k keeper.Keeper) IBCModule {
	return IBCModule{
		keeper: k,
	}
}

// OnChanOpenInit implements the IBCModule interface
func (im IBCModule) OnChanOpenInit(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID string,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	version string,
) error {
	return sdkerrors.Wrap(icatypes.ErrInvalidChannelFlow, "channel handshake must be initiated by controller chain")
}

// OnChanOpenTry implements the IBCModule interface
func (im IBCModule) OnChanOpenTry(
	ctx sdk.Context,
	order channeltypes.Order,
	connectionHops []string,
	portID,
	channelID string,
	chanCap *capabilitytypes.Capability,
	counterparty channeltypes.Counterparty,
	counterpartyVersion string,
) (string, error) {
	if !im.keeper.IsHostEnabled(ctx) {
		return "", types.ErrHostSubModuleDisabled
	}

	return im.keeper.OnChanOpenTry(ctx, order, connectionHops, portID, channelID, chanCap, counterparty, counterpartyVersion)
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenAck(
	ctx sdk.Context,
	portID,
	channelID string,
	counterpartyChannelID string,
	counterpartyVersion string,
) error {
	return sdkerrors.Wrap(icatypes.ErrInvalidChannelFlow, "channel handshake must be initiated by controller chain")
}

// OnChanOpenAck implements the IBCModule interface
func (im IBCModule) OnChanOpenConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	if !im.keeper.IsHostEnabled(ctx) {
		return types.ErrHostSubModuleDisabled
	}

	return im.keeper.OnChanOpenConfirm(ctx, portID, channelID)
}

// OnChanCloseInit implements the IBCModule interface
func (im IBCModule) OnChanCloseInit(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	// Disallow user-initiated channel closing for interchain account channels
	return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "user cannot close channel")
}

// OnChanCloseConfirm implements the IBCModule interface
func (im IBCModule) OnChanCloseConfirm(
	ctx sdk.Context,
	portID,
	channelID string,
) error {
	return im.keeper.OnChanCloseConfirm(ctx, portID, channelID)
}

// OnRecvPacket implements the IBCModule interface
func (im IBCModule) OnRecvPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	_ sdk.AccAddress,
) ibcexported.Acknowledgement {
	if !im.keeper.IsHostEnabled(ctx) {
		return types.NewErrorAcknowledgement(types.ErrHostSubModuleDisabled)
	}

	txResponse, err := im.keeper.OnRecvPacket(ctx, packet)
	if err != nil {
		// Emit an event including the error msg
		keeper.EmitWriteErrorAcknowledgementEvent(ctx, packet, err)

		return types.NewErrorAcknowledgement(err)
	}

	// NOTE: acknowledgement will be written synchronously during IBC handler execution.
	return channeltypes.NewResultAcknowledgement(txResponse)
}

// OnAcknowledgementPacket implements the IBCModule interface
func (im IBCModule) OnAcknowledgementPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	acknowledgement []byte,
	relayer sdk.AccAddress,
) error {
	return sdkerrors.Wrap(icatypes.ErrInvalidChannelFlow, "cannot receive acknowledgement on a host channel end, a host chain does not send a packet over the channel")
}

// OnTimeoutPacket implements the IBCModule interface
func (im IBCModule) OnTimeoutPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	relayer sdk.AccAddress,
) error {
	return sdkerrors.Wrap(icatypes.ErrInvalidChannelFlow, "cannot cause a packet timeout on a host channel end, a host chain does not send a packet over the channel")
}
