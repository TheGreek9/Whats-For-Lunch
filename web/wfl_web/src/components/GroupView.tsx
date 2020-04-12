import * as React from "react"
import { useState, useEffect } from "react"
import Typography from "@material-ui/core/Typography"
import axios from "axios"

import { RouteProps } from "../App"

export const GroupView = (props: RouteProps) => {
    const { match } = props
    const [testData, setTestData] = useState({groupID: 0, groupName: ""})
    console.log(match)
    useEffect(() => {
        axios({
            url: `http://localhost:8080/v1/group/query/${match.params.groupId}`,
            method: "get",
        })
        .then(res => {
            const data = res.data
            console.log(data)
            setTestData(data)
        })
    }, [match.params.groupId])

    return <Typography variant="h5">
        Group ID: {testData.groupID}
        <br/>
        Group Name: {testData.groupName}
        </Typography>
}