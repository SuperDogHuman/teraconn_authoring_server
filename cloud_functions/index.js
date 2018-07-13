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

  const speech  = require('@google-cloud/speech');
  const client  = new speech.v1.SpeechClient();

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
        fileID:   fileID,
        lessonID: lessonID,
        text:     transcription,
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

  const fileID = file.name.split('-')[1].slice(0, -4);
  const info = {
    fileID:      fileID,
    isConverted: true,
  };

  recordVoiceInfo(fileID, info)
    .then(() => { callback(); })
    .catch((err) => {
      console.error(err);
      callback();
    });
}

function recordVoiceInfo(fileID, info) {
  const datastore   = require('@google-cloud/datastore')();
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
      console.error(err);
      return transaction.rollback();
    });
}
