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
			    <ProfileImageComponent :imageUrl="profile.ProfilePictureLocation == ''? '/images/icon_standard.png': profile.ProfilePictureLocation" ></ProfileImageComponent>
                {{profile.Username}}
            </RouterLink>
        </div>
    </div>
</template>

<style>


</style>
