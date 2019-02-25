<template>
  <div class="contact-item">
    <p>
      <span class="item-title">{{contact.name}}</span>
      <span class="item-count">({{contact.msgCount}})</span>
    </p>
    <span class="item-time"><timeago :datetime="contact.time" :auto-update="60"></timeago></span>
  </div>
</template>

<style>
.contact-item {
  font-size: 13px;
  padding: 10px 10px;
  border-bottom: 1px solid #ddd;
  cursor: pointer;
}

.contact-item:hover {
  background: rgba(0, 123, 255, .1);
  border-bottom: 1px solid rgba(0, 123, 255, .2);
}

.contact-item p {
  margin: 0;
  padding: 0;
  font-weight: bold;
  color: #333;
}

.contact-item .item-time {
  font-size: 11px;
  color: #075419;
}

</style>

<script>
export default {
  props: ['contact']
}
</script>
