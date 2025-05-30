/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/

import * as fm from "../../fetch.pb"
import * as ZjutJhApiTypesV1alpha1Response from "../types/v1alpha1/response.pb"
import * as ZjutJhApiUserV1alpha1User from "../user/v1alpha1/user.pb"

export class User {
    static Register(req: ZjutJhApiUserV1alpha1User.RegisterRequest, initReq?: fm.InitReq): Promise<ZjutJhApiTypesV1alpha1Response.Response> {
        return fm.fetchReq<ZjutJhApiUserV1alpha1User.RegisterRequest, ZjutJhApiTypesV1alpha1Response.Response>(`/api/register`, {
            ...initReq,
            method: "POST",
            body: JSON.stringify(req, fm.replacer)
        })
    }

    static Login(req: ZjutJhApiUserV1alpha1User.LoginRequest, initReq?: fm.InitReq): Promise<ZjutJhApiTypesV1alpha1Response.Response> {
        return fm.fetchReq<ZjutJhApiUserV1alpha1User.LoginRequest, ZjutJhApiTypesV1alpha1Response.Response>(`/api/auth`, {
            ...initReq,
            method: "POST",
            body: JSON.stringify(req, fm.replacer)
        })
    }

    static ResetPassword(req: ZjutJhApiUserV1alpha1User.ResetPasswordRequest, initReq?: fm.InitReq): Promise<ZjutJhApiTypesV1alpha1Response.Response> {
        return fm.fetchReq<ZjutJhApiUserV1alpha1User.ResetPasswordRequest, ZjutJhApiTypesV1alpha1Response.Response>(`/api/repass`, {
            ...initReq,
            method: "POST",
            body: JSON.stringify(req, fm.replacer)
        })
    }

    static Delete(req: ZjutJhApiUserV1alpha1User.DeleteRequest, initReq?: fm.InitReq): Promise<ZjutJhApiTypesV1alpha1Response.Response> {
        return fm.fetchReq<ZjutJhApiUserV1alpha1User.DeleteRequest, ZjutJhApiTypesV1alpha1Response.Response>(`/api/del`, {
            ...initReq,
            method: "POST",
            body: JSON.stringify(req, fm.replacer)
        })
    }

    static OauthLogin(req: ZjutJhApiUserV1alpha1User.LoginRequest, initReq?: fm.InitReq): Promise<ZjutJhApiTypesV1alpha1Response.Response> {
        return fm.fetchReq<ZjutJhApiUserV1alpha1User.LoginRequest, ZjutJhApiTypesV1alpha1Response.Response>(`/api/oauth`, {
            ...initReq,
            method: "POST",
            body: JSON.stringify(req, fm.replacer)
        })
    }
}