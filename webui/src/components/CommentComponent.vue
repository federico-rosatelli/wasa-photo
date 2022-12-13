<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
            // idUser:null,
            // follow:null,
        }
    },
	props: ["commentData"],
}
</script>

<template>
    <table id="table">
        <tr>
            <td>Username</td>
            <td>Comment</td>
            <td>Time</td>
        </tr>
        <tr v-for="comment in commentData">
            <td>
                <RouterLink v-bind:to="'/profile/'+comment.Username" class="nav-link">
                    <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                    {{comment.Username}}
                </RouterLink>
            </td>
            <td>
                {{comment.Content}}
            </td>
            <td>
                {{comment.Time}}
            </td>
        </tr>
    </table>
    
</template>

<style>
.row{
 display: flex;
}
td{
    font-size: 15px;
}
#table {
 font-family: Arial, Helvetica, sans-serif;
 border-collapse: collapse;
 width: 100%;
}
#table td, #table th {
 border: 1px solid #ddd;
 padding: 20px;
}
#table .title{
  font-size: 25px;
  background-color: #999999;
  text-align: center;
}
#table tr:nth-child(even){background-color: #f2f2f2;}
#table tr:hover {background-color: #ddd;}
#table th {
 padding-top: 12px;
 padding-bottom: 12px;
 text-align: left;
 background-color: #04AA6D;
 color: white;
}

#table a{
  color: #000;
}
</style>