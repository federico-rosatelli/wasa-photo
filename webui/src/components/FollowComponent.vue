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
	props: ["usersFollow"],
}
</script>

<template>
	<div v-for="profile in usersFollow">
        <div style="display: flex;">
            <h2>{{profile.Username}}</h2>
            <h2>{{profile.Id}}</h2>
        </div>
    </div>
</template>

<style>


</style>
