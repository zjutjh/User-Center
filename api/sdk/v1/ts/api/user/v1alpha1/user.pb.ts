/* eslint-disable */
// @ts-nocheck
/*
* This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
*/
export type RegisterRequest = {
    studentId?: string
    password?: string
    iid?: string
    email?: string
}

export type LoginRequest = {
    studentId?: string
    password?: string
}

export type ResetPasswordRequest = {
    studentId?: string
    password?: string
    iid?: string
}

export type DeleteRequest = {
    studentId?: string
    iid?: string
}