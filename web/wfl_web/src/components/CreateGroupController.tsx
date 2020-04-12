import * as React from "react"
import axios from "axios"
import { Method } from "axios"
import { useState } from "react"
import TextField from "@material-ui/core/TextField"
import Button from "@material-ui/core/Button"

import { RouteProps } from "../App"

export interface Group {
    groupID?: number
	groupName?: string
	groupMembers?: number[]
	restaurantsTried?: number[]
	restaurantsMissed?: number[]
	collectiveZip?: string

}

export const CreateGroupController = (props: RouteProps) => {
    const [groupName, setGroupName] = useState("")
    const [passFail, setPassFail] = useState<"success" | "failure">("failure")
    const onClick = () => {
        console.log("on Click")
        const req: {url: string, method: Method, data: Group} = {
            url: `http://localhost:8080/v1/group/create`,
            method: "post",
            data: {
                groupName: groupName,
                groupMembers: [1, 2, 3],
                restaurantsTried: [1, 2, 3],
                restaurantsMissed: [1, 2, 3],
                collectiveZip: "01234",
            }
        }
        axios(req).then(resp => {
            console.log(resp.data)
            const data: Group = resp.data
            setPassFail("success")
        })
    }

    return (
        <>
        <TextField onChange={ev => setGroupName(ev.target.value)}/>
        <Button onClick={onClick}>Submit</Button>
        <p style={{color: passFail === "failure" ? "red" : "green"}}>{passFail}</p>
        </>
    )
}