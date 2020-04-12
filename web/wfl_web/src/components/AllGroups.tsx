import * as React from "react"
import { useState, useEffect } from "react"

export const AllGroups = () => {
    return
}

// import { Group } from "./CreateGroupController"

// export const AllGroups = () => {
//     const [currentGroups, setCurrentGroups] = useState<Group[]>([])

//     useEffect(() => {
//         axios({

//         })
//     })

    
//     return <TableContainer component={Paper}>
//     <Table className={classes.table} aria-label="simple table">
//       <TableHead>
//         <TableRow>
//           <TableCell>Dessert (100g serving)</TableCell>
//           <TableCell align="right">Calories</TableCell>
//           <TableCell align="right">Fat&nbsp;(g)</TableCell>
//           <TableCell align="right">Carbs&nbsp;(g)</TableCell>
//           <TableCell align="right">Protein&nbsp;(g)</TableCell>
//         </TableRow>
//       </TableHead>
//       <TableBody>
//         {rows.map((row) => (
//           <TableRow key={row.name}>
//             <TableCell component="th" scope="row">
//               {row.name}
//             </TableCell>
//             <TableCell align="right">{row.calories}</TableCell>
//             <TableCell align="right">{row.fat}</TableCell>
//             <TableCell align="right">{row.carbs}</TableCell>
//             <TableCell align="right">{row.protein}</TableCell>
//           </TableRow>
//         ))}
//       </TableBody>
//     </Table>
//   </TableContainer>
// }