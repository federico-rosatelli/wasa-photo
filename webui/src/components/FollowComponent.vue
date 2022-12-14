<script>
export default {
    data:()=>{
        return {
            profiles: null,
            errormsg:null,
            show:false,
            // idUser:null,
            // follow:null,
        }
    },
	props: ["usersFollow","showModal"],
}
</script>

<template>
    <div v-for="profile in usersFollow">
        <div style="display: flex;">
            <RouterLink v-bind:to="'/profile/'+profile.Username" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
                {{profile.Username}}
            </RouterLink>
            <h2>{{profile.Id}}</h2>
        </div>
    </div>
</template>

<style>


</style>
