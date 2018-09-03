const ffmpeg     = require('fluent-ffmpeg');
const ffmpegPath = require('ffmpeg-static').path;

const storage    = require('@google-cloud/storage')();
const datastore  = require('@google-cloud/datastore')();

exports.wavToText = async (data, context) => {
    const file = data;
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
    const responses = await client.recognize(request).catch((err) => {
        console.error(err);
    });

    const response = responses[0];
    const transcription = response.results
        .map((result) => result.alternatives[0].transcript)
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
    recordVoiceInfo(fileID, info);
}

exports.wavTo16Khz = async (data, context) => {
    const file = data;
    if (file.size == 0) { return; }

    const fileNames                 = file.name.split('-');
    const fileID                    = fileNames[1].slice(0, -4);
    const localRawWAVFilePath       = `/tmp/raw-${fileID}.wav`;
    const localResampledWAVFilePath = `/tmp/resampled-${fileID}.wav`;
    const remoteWAVFilePath         = file.name;

    await downloadFromStorage('teraconn_raw_voice', file.name, localRawWAVFilePath).catch((err) => {
        console.error(err);
    });

    await encodeToWav(localRawWAVFilePath, localResampledWAVFilePath).catch((err) => {
        console.error(err);
    });

    await uploadToStorage('teraconn_voice_for_transcription', localResampledWAVFilePath, remoteWAVFilePath).catch((err) => {
        console.error(err);
    });

    function encodeToWav(rawFilePath, resampledFilePath) {
        return new Promise((resolve, reject) => {
            ffmpeg(rawFilePath)
                .setFfmpegPath(ffmpegPath)
                .audioChannels(1)
                .audioFrequency(16000)
                .output(resampledFilePath)
                .on('end', () => { resolve(); })
                .on('error', (err) => {
                    console.error(err);
                    reject(error);
                })
                .run();
        });
    }
}

exports.wavToOgg = async (data, context) => {
    const file = data;
    if (file.size == 0) { return; }

    const fileNames         = file.name.split('-');
    const lessonID          = fileNames[0];
    const fileID            = fileNames[1].slice(0, -4);
    const localWAVFilePath  = `/tmp/${fileID}.wav`;
    const localOGGFilePath  = `/tmp/${fileID}.ogg`;
    const remoteOGGFilePath = `voice/${lessonID}/${fileID}.ogg`;

    await downloadFromStorage('teraconn_raw_voice', file.name, localWAVFilePath).catch((err) => {
        console.error(err);
    });

    await encodeToOgg(localWAVFilePath, localOGGFilePath).catch((err) => {
        console.error(err);
    });

    await uploadToStorage('teraconn_material', localOGGFilePath, remoteOGGFilePath).catch((err) => {
        console.error(err);
    });

    const info = {
        LessonID:    lessonID,
        FileID:      fileID,
        IsConverted: true,
    };
    await recordVoiceInfo(fileID, info).catch((err) => {
        console.error(err);
    });

    function encodeToOgg(wavFilePath, oggFilePath) {
        return new Promise((resolve, reject) => {
            ffmpeg(wavFilePath)
                .setFfmpegPath(ffmpegPath)
                .audioCodec('libvorbis')
                .audioQuality(0)
                .audioChannels(1)
//                .audioFilters(
//                    [{ filter: 'dynaudnorm', options: '' }]
//                )
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
}

function downloadFromStorage(bucketName, remoteFilePath, localFilePath) {
    return storage.bucket(bucketName)
        .file(remoteFilePath)
        .download({ destination: localFilePath });
}

function uploadToStorage(bucketName, localFilePath, remoteFilePath) {
    const bucket = storage.bucket(bucketName);
    const options = {
        destination: bucket.file(remoteFilePath),
        resumable:   false,
    };
    return bucket.upload(localFilePath, options);
}

function recordVoiceInfo(fileID, info) {
    const transaction = datastore.transaction();
    const key         = datastore.key(['VoiceText', fileID]);
    return transaction.run()
        .then((data) => {
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