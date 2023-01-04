<script>
export default {
	data: function() {
		return {
			errormsg: null,
			loading: false,
			some_data: null,
			token: null,
			img_src: null,
		}
	},
	methods: {
		async sendPhoto(){
			var formData = new FormData();
			formData.append('text',document.getElementById("text").value)
			formData.append('myFile',document.getElementById("myFile").files[0])
			try {
				let response = await this.$axios.post("/addphoto",formData,{headers:{"Token":this.token}});
				this.some_data = response.data;
			} catch (e) {
				this.errormsg = e.toString();
			}
			this.loading = false;
			this.$router.push("profile")
		},
		onFileChange(e){
			// let imgPrev = document.getElementById("previewImg")
			const file = e.target.files[0];
			// imgPrev.src = URL.createObjectURL(file)
			this.img_src = URL.createObjectURL(file)
		},
	},
}
</script>

<template>
	<div>
		<div
			class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
			<h1 class="h2">Upload Image</h1>
			<div class="btn-toolbar mb-2 mb-md-0">
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="refresh">
						Refresh
					</button>
					<button type="button" class="btn btn-sm btn-outline-secondary" @click="exportList">
						Export
					</button>
				</div>
				<div class="btn-group me-2">
					<button type="button" class="btn btn-sm btn-outline-primary" @click="newItem">
						New
					</button>
				</div>
			</div>
		</div>

		<div>
			<input type="file" id="myFile" accept=".png, .jpg, .jpeg" @change="onFileChange">
			<input type="text" placeholder="Description" id="text">
			<button type="button" class="btn btn-sm btn-outline-primary" @click="sendPhoto">Send</button>
		</div>
		<img class="box" v-if="img_src" v-bind:src="img_src" alt="Your Image" id="previewImg">

		<ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
	</div>
</template>
<style>
</style>
