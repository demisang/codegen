import axios from 'axios'
import Placeholder from "@/models/placeholder";

// const baseUrl = process.env.BASE_URL
const baseUrl = "http://localhost:4765"

type onSuccessFunction = (data: any, statusCode: number) => void;
type onFailureFunction = (data: any, statusCode: number, httpStatusMessage: string) => void;
type onStartFunction = () => void;
type onCompletedFunction = () => void;

function sendRequest(
    url: string,
    method: string,
    params: object | null,
    onSuccess: onSuccessFunction,
    onFailure: onFailureFunction,
    onStart: onStartFunction,
    onCompleted: onCompletedFunction,
) {
    if (method === "GET" && params !== null) {
        const queryParams = new URLSearchParams()
        for (const paramsKey in params) {
            // queryParams.set(paramsKey, params[paramsKey])
        }

        url += "?" + queryParams.toString()
    }

    fetch(baseUrl + url, {
        method: method,
        headers: {
            "Content-Type": "application/json",
        },
        body: method === "GET" && params !== null ? JSON.stringify(params) : null,
    })
}

export function getTemplates(
    type = "",
    onSuccess: onSuccessFunction = () => null,
    onFailure: onFailureFunction = () => null,
    onStart: onStartFunction = () => null,
    onCompleted: onCompletedFunction = () => null,
) {
    onStart()
    axios.get(baseUrl + "/templates", {params: {type: type}})
        .then((response) => {
            onSuccess(response.data, response.status)
        })
        .catch((error) => {
            onFailure(error.data, error.statusCode, error.statusText)
        })
        .finally(() => {
            onCompleted()
        })
}

export function getRawList(
    templateID: string,
    targetDir: string,
    placeholders: Placeholder[],
    onSuccess: onSuccessFunction = () => null,
    onFailure: onFailureFunction = () => null,
    onStart: onStartFunction = () => null,
    onCompleted: onCompletedFunction = () => null,
) {
    onStart()
    axios.post(baseUrl + "/raw-list", {template_id: templateID, target_dir: targetDir, placeholders: placeholders})
        .then((response) => {
            onSuccess(response.data, response.status)
        })
        .catch((error) => {
            onFailure(error.data, error.statusCode, error.statusText)
        })
        .finally(() => {
            onCompleted()
        })
}

export function getDirectories(
    selected = "",
    onSuccess: onSuccessFunction = () => null,
    onFailure: onFailureFunction = () => null,
    onStart: onStartFunction = () => null,
    onCompleted: onCompletedFunction = () => null,
) {
    onStart()
    axios.get(baseUrl + "/directories", {params: {selected: selected}})
        .then((response) => {
            onSuccess(response.data, response.status)
        })
        .catch((error) => {
            onFailure(error.data, error.statusCode, error.statusText)
        })
        .finally(() => {
            onCompleted()
        })
}

export function generate(
    templateID: string,
    targetDir: string,
    placeholders: Placeholder[],
    onSuccess: onSuccessFunction = () => null,
    onFailure: onFailureFunction = () => null,
    onStart: onStartFunction = () => null,
    onCompleted: onCompletedFunction = () => null,
) {
    onStart()
    const data = {
        template_id: templateID,
        target_dir: targetDir,
        placeholders: placeholders,
    }
    axios.post(baseUrl + "/generate", data)
        .then((response) => {
            onSuccess(response.data, response.status)
        })
        .catch((error) => {
            onFailure(error.data, error.statusCode, error.statusText)
        })
        .finally(() => {
            onCompleted()
        })
}
