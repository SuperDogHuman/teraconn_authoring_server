const ffmpeg     = require('fluent-ffmpeg');
const ffmpegPath = require('ffmpeg-static').path;

const storage    = require('@google-cloud/storage')();
const datastore  = require('@google-cloud/datastore')();

exports.wavToText = (event, callback) => {
  const file = event.data;

  if (file.size == 0) { return callback(); }

  const request = {
    audio: {
      uri: `gs://${file.bucket}/${file.name}`
    },
    config: {
      encoding:        'LINEAR16',
      sampleRateHertz: 16000,
      languageCode:    'ja-JP',
    },
  };

  const speech = require('@google-cloud/speech');
  const client = new speech.v1.SpeechClient();

  client.recognize(request)
    .then(data => {
      const response = data[0];
      const transcription = response.results
        .map(result => result.alternatives[0].transcript)
        .join('\n');

      const fileNames = file.name.split('-');
      const lessonID  = fileNames[0];
      const fileID    = fileNames[1].slice(0, -4);
      const info = {
        LessonID: lessonID,
        FileID:   fileID,
        IsTexted: true,
        Text:     transcription,
      };
      recordVoiceInfo(fileID, info)
    })
    .then(() => { callback(); })
    .catch((err) => {
      console.error(err);
      callback();
    });
}

exports.wavToOgg = (event, callback) => {
  const file = event.data;
  if (file.size == 0) { return callback(); }

  const fileNames         = file.name.split('-');
  const lessonID          = fileNames[0];
  const fileID            = fileNames[1].slice(0, -4);
  const localWAVFilePath  = `/tmp/${fileID}.wav`;
  const localOGGFilePath  = `/tmp/${fileID}.ogg`;
  const remoteOGGFilePath = `voice/${lessonID}/${fileID}.ogg`;

  downloadFromStorage(file.name, localWAVFilePath)
    .then(() => {
      return encodeToOgg(localWAVFilePath, localOGGFilePath);
    })
    .then(() => {
      return uploadToStorage(localOGGFilePath, remoteOGGFilePath);
    })
    .then(() => {
      const info = {
        LessonID:    lessonID,
        FileID:      fileID,
        IsConverted: true,
      };
      return recordVoiceInfo(fileID, info);
    })
    .then(() => { callback(); })
    .catch((err) => {
      console.error(err);
      callback();
    });

  function downloadFromStorage(remoteFilePath, localFilePath) {
    return storage.bucket('teraconn_raw_voice')
      .file(remoteFilePath)
      .download({ destination: localFilePath });
  }

  function encodeToOgg(wavFilePath, oggFilePath) {
    return new Promise((resolve, reject) => {
      ffmpeg(wavFilePath)
        .setFfmpegPath(ffmpegPath)
        .audioChannels(1)
        .audioFrequency(16000)
        .format('ogg')
        .output(oggFilePath)
        .on('end', () => { resolve(); })
        .on('error', (err) => {
          console.error(err);
          reject(error);
        })
        .run();
    });
  }

  function uploadToStorage(localFilePath, remoteFilePath) {
    const bucket = storage.bucket('teraconn_material');
    const options = {
      destination: bucket.file(remoteFilePath),
      resumable:   false,
    };
    return bucket.upload(localFilePath, options);
  }
}

function recordVoiceInfo(fileID, info) {
  const transaction = datastore.transaction();
  const key         = datastore.key(['VoiceText', fileID]);
  return transaction.run()
    .then(function(data) {
      return transaction.get(key);
    })
    .then((results) => {
      let entity = results[0];
      if (entity) {
        Object.keys(info).forEach((k) => {
          entity[k] = info[k];
        });
        return transaction.save({ key: key, data: entity });
      } else {
        entity = { key: key, data: info};
        return transaction.save(entity);
      }
    })
    .then(() =>{
      return transaction.commit();
    })
    .catch((err) => {
      if (err.code == 10) {
        console.warn("too much contention, retry after 1 second");
        return setTimeout(recordVoiceInfo.bind(this, fileID, info), 1000);
      } else {
        console.error(err);
        return transaction.rollback();
      }
    });
}