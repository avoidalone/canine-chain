// THIS FILE IS GENERATED AUTOMATICALLY. DO NOT MODIFY.

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient } from "@cosmjs/stargate";
import { Registry, OfflineSigner, EncodeObject, DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import { Api } from "./rest";
import { MsgPostkey } from "./types/filetree/tx";
import { MsgAddViewers } from "./types/filetree/tx";
import { MsgDeleteFile } from "./types/filetree/tx";
import { MsgPostFile } from "./types/filetree/tx";
import { MsgInitAll } from "./types/filetree/tx";
import { MsgInitAccount } from "./types/filetree/tx";


const types = [
  ["/jackaldao.canine.filetree.MsgPostkey", MsgPostkey],
  ["/jackaldao.canine.filetree.MsgAddViewers", MsgAddViewers],
  ["/jackaldao.canine.filetree.MsgDeleteFile", MsgDeleteFile],
  ["/jackaldao.canine.filetree.MsgPostFile", MsgPostFile],
  ["/jackaldao.canine.filetree.MsgInitAll", MsgInitAll],
  ["/jackaldao.canine.filetree.MsgInitAccount", MsgInitAccount],
  
];
export const MissingWalletError = new Error("wallet is required");

export const registry = new Registry(<any>types);

const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
}

interface SignAndBroadcastOptions {
  fee: StdFee,
  memo?: string
}

const txClient = async (wallet: OfflineSigner, { addr: addr }: TxClientOptions = { addr: "http://localhost:26657" }) => {
  if (!wallet) throw MissingWalletError;
  let client;
  if (addr) {
    client = await SigningStargateClient.connectWithSigner(addr, wallet, { registry });
  }else{
    client = await SigningStargateClient.offline( wallet, { registry });
  }
  const { address } = (await wallet.getAccounts())[0];

  return {
    signAndBroadcast: (msgs: EncodeObject[], { fee, memo }: SignAndBroadcastOptions = {fee: defaultFee, memo: ""}) => client.signAndBroadcast(address, msgs, fee,memo),
    msgPostkey: (data: MsgPostkey): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgPostkey", value: MsgPostkey.fromPartial( data ) }),
    msgAddViewers: (data: MsgAddViewers): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgAddViewers", value: MsgAddViewers.fromPartial( data ) }),
    msgDeleteFile: (data: MsgDeleteFile): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgDeleteFile", value: MsgDeleteFile.fromPartial( data ) }),
    msgPostFile: (data: MsgPostFile): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgPostFile", value: MsgPostFile.fromPartial( data ) }),
    msgInitAll: (data: MsgInitAll): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgInitAll", value: MsgInitAll.fromPartial( data ) }),
    msgInitAccount: (data: MsgInitAccount): EncodeObject => ({ typeUrl: "/jackaldao.canine.filetree.MsgInitAccount", value: MsgInitAccount.fromPartial( data ) }),
    
  };
};

interface QueryClientOptions {
  addr: string
}

const queryClient = async ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseUrl: addr });
};

export {
  txClient,
  queryClient,
};
