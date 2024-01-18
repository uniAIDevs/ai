// @generated by protoc-gen-connect-es v1.2.0 with parameter "target=ts,import_extension="
// @generated from file llm/service.proto (package llm.v1beta1, syntax proto3)
/* eslint-disable */
// @ts-nocheck

import { ChatPrediction, ChatRequest, CompletionPrediction, CompletionRequest, LoadModelRequest, LoadModelResponse } from "./service_pb";
import { MethodKind } from "@bufbuild/protobuf";

/**
 * @generated from service llm.v1beta1.LLMService
 */
export const LLMService = {
  typeName: "llm.v1beta1.LLMService",
  methods: {
    /**
     * @generated from rpc llm.v1beta1.LLMService.Chat
     */
    chat: {
      name: "Chat",
      I: ChatRequest,
      O: ChatPrediction,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc llm.v1beta1.LLMService.Completion
     */
    completion: {
      name: "Completion",
      I: CompletionRequest,
      O: CompletionPrediction,
      kind: MethodKind.Unary,
    },
    /**
     * @generated from rpc llm.v1beta1.LLMService.LoadModel
     */
    loadModel: {
      name: "LoadModel",
      I: LoadModelRequest,
      O: LoadModelResponse,
      kind: MethodKind.Unary,
    },
  }
} as const;
