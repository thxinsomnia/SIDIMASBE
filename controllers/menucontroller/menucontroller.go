package menucontroller

import (

    "fmt"
    "strconv"
	"net/http"
	"SIDIMASBE/models"

	"github.com/gin-gonic/gin"
)

func BuatMenu(c *gin.Context) {
    var menuInput struct {
        IDMenu      int64                  `json:"id_menu"`
        NamaMenu    string                 `json:"nama_menu"`
        Deskripsi   string                 `json:"deskripsi"`
        JumlahPorsi int64                   `json:"jumlah_porsi"`
        Bahan       []models.MenuBahanItem `json:"bahan"`
    }

    if err := c.ShouldBindJSON(&menuInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    tx := models.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Terjadi kesalahan internal"})
        }
    }()

    // Jika ID tidak diberikan, generate otomatis
    if menuInput.IDMenu == 0 {
        newID, err := GenerateMenuID()
        if err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal membuat ID menu"})
            return
        }
        menuInput.IDMenu, _ = strconv.ParseInt(newID, 10, 64) // Konversi ke int64
    }

    // Pastikan ID unik
    var existingMenu models.Menu
    if err := tx.First(&existingMenu, "id_menu = ?", menuInput.IDMenu).Error; err == nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"Message": "ID menu sudah digunakan, coba lagi"})
        return
    }

    // Buat menu baru
    menu := models.Menu{
        ID_menu:      menuInput.IDMenu,
        Nama_menu:    menuInput.NamaMenu,
        Deskripsi:    menuInput.Deskripsi,
        Jumlah_porsi: menuInput.JumlahPorsi,
    }
    if err := tx.Create(&menu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal membuat menu"})
        return
    }

    // Tambahkan bahan-bahan ke menu
    for _, item := range menuInput.Bahan {
        var bahan models.Material
        if err := tx.First(&bahan, item.ID_bahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf("Bahan dengan ID %d tidak ditemukan", item.ID_bahan)})
            return
        }

        menuBahan := models.MenuBahan{
            ID_menu:   menu.ID_menu,
            ID_bahan:  item.ID_bahan,
            Kebutuhan: item.Kebutuhan,
        }
        if err := tx.Create(&menuBahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal menambahkan bahan ke menu"})
            return
        }
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Menu berhasil dibuat", "id_menu": menu.ID_menu})
}


func GenerateMenuID() (string, error) {
    var lastMenuID string

    // Ambil ID terakhir dari database
    err := models.DB.Model(&models.Menu{}).
        Select("id_menu").
        Order("id_menu DESC").
        Limit(1).
        Pluck("id_menu", &lastMenuID).
        Error

    if err != nil || lastMenuID == "" {
        // Jika tidak ada data sebelumnya, mulai dari 3399001
        return "7142001", nil
    }

    // Pastikan ID terakhir memiliki panjang yang sesuai
    if len(lastMenuID) < 7 {
        return "7142001", nil
    }

    // Ambil bagian angka urut (3 digit terakhir)
    lastIDNumeric, err := strconv.Atoi(lastMenuID[4:]) // Ambil setelah "3399"
    if err != nil {
        return "", err
    }

    // Tambahkan 1 ke ID terakhir
    newIDNumeric := lastIDNumeric + 1
    newID := fmt.Sprintf("7142%03d", newIDNumeric) // Format: "3399XXX"

    return newID, nil
}

func AmbilDataMenu(c *gin.Context) {
	var menus []models.Menu

	if err := models.DB.Preload("Bahan.Materials.Supplier").Find(&menus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data menu"})
		return
	}

	// Struct untuk menampilkan hanya informasi yang dibutuhkan
	type BahanResponse struct {
		NamaBahan    string `json:"nama_bahan"`
		NamaSupplier string `json:"nama_supplier"`
        Kebutuhan int64 `json:"kebutuhan"`
	}

	type MenuResponse struct {
		IDMenu      int64           `json:"id_menu"`
		NamaMenu    string          `json:"nama_menu"`
		Deskripsi   string          `json:"deskripsi"`
		JumlahPorsi int64           `json:"jumlah_porsi"`
		Bahan       []BahanResponse `json:"bahan"`
	}

	var response []MenuResponse

	// Looping untuk memformat data sebelum dikirim sebagai respons
	for _, menu := range menus {
		var bahanList []BahanResponse
		for _, menuBahan := range menu.Bahan {
			bahanList = append(bahanList, BahanResponse{
				NamaBahan:    menuBahan.Materials.Nama_bahan,
				NamaSupplier: menuBahan.Materials.Supplier.Nama_supplier,
                Kebutuhan: menuBahan.Kebutuhan,
			})
		}

		response = append(response, MenuResponse{
			IDMenu:      menu.ID_menu,
			NamaMenu:    menu.Nama_menu,
			Deskripsi:   menu.Deskripsi,
			JumlahPorsi: menu.Jumlah_porsi,
			Bahan:       bahanList,
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": response})
}


func AmbilDataMenuID(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu
	if err := models.DB.Preload("Bahan").First(&menu, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Menu tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": menu})
}

func EditMenu(c *gin.Context) {
    id := c.Param("id")

    var menuInput struct {
        NamaMenu    string                  `json:"nama_menu"`
        Deskripsi   string                  `json:"deskripsi"`
        JumlahPorsi int64                   `json:"jumlah_porsi"`
        Bahan       []models.MenuBahanItem  `json:"bahan"`
    }

    if err := c.ShouldBindJSON(&menuInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    tx := models.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Terjadi kesalahan internal"})
        }
    }()

    var menu models.Menu
    if err := tx.First(&menu, id).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"Message": "Menu tidak ditemukan"})
        return
    }

    // Update menu
    menu.Nama_menu = menuInput.NamaMenu
    menu.Deskripsi = menuInput.Deskripsi
    menu.Jumlah_porsi = menuInput.JumlahPorsi
    if err := tx.Save(&menu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengupdate menu"})
        return
    }

    // Hapus bahan lama
    if err := tx.Where("id_menu = ?", id).Delete(&models.MenuBahan{}).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal menghapus bahan lama"})
        return
    }

    // Tambahkan bahan baru dengan validasi
    for _, item := range menuInput.Bahan {
        var bahan models.Material
        if err := tx.First(&bahan, item.ID_bahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{"Message": fmt.Sprintf("Bahan dengan ID %d tidak ditemukan", item.ID_bahan)})
            return
        }

        menuBahan := models.MenuBahan{
            ID_menu:   menu.ID_menu,
            ID_bahan:  item.ID_bahan,
            Kebutuhan: item.Kebutuhan,
        }
        if err := tx.Create(&menuBahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal menambahkan bahan baru"})
            return
        }
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"Message": "Menu berhasil diupdate"})
}


func HapusMenu(c *gin.Context) {
	id := c.Param("id")
	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("id_menu = ?", id).Delete(&models.MenuBahan{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus bahan terkait"})
		return
	}

	if err := tx.Delete(&models.Menu{}, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus menu"})
		return
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Menu berhasil dihapus"})
}

func BuatPorsiMenu(c *gin.Context) {
    var pesananInput struct {
        IDMenu      int64 `json:"id_menu"`
        JumlahPorsi int64 `json:"jumlah_porsi"`
    }

    if err := c.ShouldBindJSON(&pesananInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
        return
    }

    tx := models.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Terjadi kesalahan internal"})
        }
    }()

    // Cek apakah menu ada
    var menu models.Menu
    if err := tx.First(&menu, "id_menu = ?", pesananInput.IDMenu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{"Message": "Menu tidak ditemukan"})
        return
    }

    // Ambil bahan-bahan yang terkait dengan menu ini
    var menuBahans []models.MenuBahan
    if err := tx.Where("id_menu = ?", pesananInput.IDMenu).Find(&menuBahans).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengambil data bahan"})
        return
    }

    bahanMap := make(map[int64]models.Material)

    // Cek apakah stok bahan cukup
    for _, menuBahan := range menuBahans {
        if _, exists := bahanMap[menuBahan.ID_bahan]; !exists {
            var bahan models.Material
            if err := tx.First(&bahan, "id_bahan = ?", menuBahan.ID_bahan).Error; err != nil {
                tx.Rollback()
                c.JSON(http.StatusInternalServerError, gin.H{"Message": fmt.Sprintf("Bahan dengan ID %d tidak ditemukan", menuBahan.ID_bahan)})
                return
            }
            bahanMap[menuBahan.ID_bahan] = bahan
        }

        totalKebutuhan := menuBahan.Kebutuhan * pesananInput.JumlahPorsi
        bahan := bahanMap[menuBahan.ID_bahan]

        if bahan.Jumlah < totalKebutuhan {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "Message": fmt.Sprintf("Stok bahan %s tidak mencukupi. Dibutuhkan %d, tetapi hanya tersedia %d.",
                    bahan.Nama_bahan, totalKebutuhan, bahan.Jumlah),
            })
            return
        }
    }

    // Kurangi stok bahan dan simpan log pengurangan bahan
    for _, menuBahan := range menuBahans {
        bahan := bahanMap[menuBahan.ID_bahan]
        totalKebutuhan := menuBahan.Kebutuhan * pesananInput.JumlahPorsi
        bahan.Jumlah -= totalKebutuhan

        if err := tx.Save(&bahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengupdate stok bahan"})
            return
        }

        // Simpan log pengurangan bahan
        logBahan := models.Log{
            ID_bahan:        bahan.ID_bahan,
            JumlahDigunakan: totalKebutuhan,
            SisaBahan:       bahan.Jumlah,
        }
        if err := tx.Create(&logBahan).Error; err != nil {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mencatat log pengurangan bahan"})
            return
        }
    }

    // UPDATE jumlah_porsi pada menu (menambah stok menu)
    menu.Jumlah_porsi += pesananInput.JumlahPorsi
    if err := tx.Save(&menu).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal mengupdate stok menu"})
        return
    }

    if err := tx.Commit().Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"Message": "Gagal commit transaction"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "Message": "Pesanan berhasil dibuat",
        "Menu":    menu.Nama_menu,
        "JumlahPorsi": pesananInput.JumlahPorsi,
        "TotalStokMenu": menu.Jumlah_porsi, // Menampilkan stok menu setelah ditambah
    })
}



