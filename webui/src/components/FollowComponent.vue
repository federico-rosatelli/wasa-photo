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
		<ProfileImageComponent :userNameF="profile.Username" :imageUrl="profile.ProfilePictureLocation == ''? '/images/icon_standard.png': profile.ProfilePictureLocation" ></ProfileImageComponent>
    </div>
</template>

<style>


</style>
